#include <alsa/asoundlib.h>
#include <cerrno>
#include <fmt123.h>
#include <mpg123.h>

#include <iostream>

//TODO - Abstract code into Classes
//     - Add Class that loads music (separate from ALSA)
//     - Add FLAC support (use FLAC API)
//     - ...

int main()
{
  // init ALSA
  snd_pcm_t* ph;
  snd_pcm_hw_params_t* hw_params;

  if (snd_pcm_open(&ph, "default", SND_PCM_STREAM_PLAYBACK, 0) < 0)
  {
    std::cout << "[ALSA] Cannot open audio device, " << stderr << std::endl;
    return -1;
  };

  if (snd_pcm_hw_params_malloc(&hw_params) < 0)
  {
    std::cout << "[ALSA] Cannot allocate hardware parameter structure, " << stderr << std::endl;
    return -1;
  }

  if (snd_pcm_hw_params_any(ph, hw_params) < 0)
  {
    std::cout << "[ALSA] Cannot initialize hardware parameter structure, " << stderr << std::endl;
    return -1;
  }


  // init MPG123
  mpg123_handle* mh;
  int err = MPG123_OK;

  mh = mpg123_new(nullptr, &err);

  if (err != 0)
    std::cout << mpg123_plain_strerror(err) << std::endl;

  long rate;
  int channels, encoding;
  if (mpg123_open(mh, "../music/yippee-tbh.mp3") || mpg123_getformat(mh, &rate, &channels, &encoding) != MPG123_OK)
    std::cout << mpg123_strerror(mh) << std::endl;
  
  // force mpg123 to use 16-bit signed if needed (since ALSA uses S16 format)
  mpg123_format(mh, rate, channels, MPG123_ENC_SIGNED_16);


  // set ALSA parameters to match MP3
  if (snd_pcm_hw_params_set_access(ph, hw_params, SND_PCM_ACCESS_RW_INTERLEAVED) < 0)
  {
    std::cout << "[ALSA] Cannot set access type, " << stderr << std::endl;
    return -1;
  }

  if (snd_pcm_hw_params_set_format(ph, hw_params, SND_PCM_FORMAT_S16_LE) < 0)
  {
    std::cout << "[ALSA] Cannot set sample format, " << stderr << std::endl;
    return -1;
  }

  if (snd_pcm_hw_params_set_rate_near(ph, hw_params, (unsigned int*)&rate, 0) < 0)
  {
    std::cout << "[ALSA] Cannot set sample rate, " << stderr << std::endl;
    return -1;
  }

  if (snd_pcm_hw_params_set_channels(ph, hw_params, channels) < 0)
  {
    std::cout << "[ALSA] Cannot set channel count, " << stderr << std::endl;
    return -1;
  }


  if (snd_pcm_hw_params(ph, hw_params) < 0)
  {
    std::cout << "[ALSA] Cannot set parameters, " << stderr << std::endl;
    return -1;
  } 

  // free the allocated memory for hardware parameters (when it is not needed anymore)
  // snd_pcm_hw_params_free(hw_params)

  if (snd_pcm_prepare(ph) < 0)
  {
    std::cout << "[ALSA] Cannot prepare audio interface for use, " << stderr << std::endl;
    return -1;
  } 

  unsigned char buffer[8192];
  size_t done;

  while (mpg123_read(mh, buffer, sizeof(buffer), &done) == MPG123_OK)
  {
    int frames = done / (channels * 2); // 2 bytes per sample since we use S16 -> 16bits = 2bytes per channel 
    int err = snd_pcm_writei(ph, (short*)buffer, frames);

    if (err == -EPIPE)
    {
      // this if statement checks whether the audio playback breaked because data wasn't delivered at time;
      // if so, continue playing the audio (prepare/reset the PCM device and try to play audio again)
      snd_pcm_prepare(ph);
    }
  }

  mpg123_close(mh);
  mpg123_delete(mh);

  // for playback wait for all pending frames to be played and then stop the PCM.
  snd_pcm_drain(ph);
  snd_pcm_close(ph);

  std::cout << "End of the Program..." << std::endl;

  return 0;
}
