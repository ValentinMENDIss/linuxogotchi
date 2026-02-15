#include "AudioSink.h"
#include "AudioDecoder.h"
#include <alsa/asoundlib.h>
#include <mpg123.h>

// here is an error somewhere! (check -> gdb ./exe -> run -> bt)
void AudioSink::PrepareAudio(AudioFile audioFile)
{
  if (snd_pcm_open(&m_PCMHandle, "default", SND_PCM_STREAM_PLAYBACK, 0) < 0)
  {
    std::cout << "[ALSA] Cannot open audio device, " << stderr << std::endl;
  };

  if (snd_pcm_hw_params_malloc(&m_HWParams) < 0)
  {
    std::cout << "[ALSA] Cannot allocate hardware parameter structure, " << stderr << std::endl;
  }

    if (snd_pcm_hw_params_any(m_PCMHandle, m_HWParams) < 0)
  {
    std::cout << "[ALSA] Cannot initialize hardware parameter structure, " << stderr << std::endl;
  }


  if (snd_pcm_hw_params_set_access(m_PCMHandle, m_HWParams, SND_PCM_ACCESS_RW_INTERLEAVED) < 0)
  {
    std::cout << "[ALSA] Cannot set access type, " << stderr << std::endl;
  }


  if (snd_pcm_hw_params_set_format(m_PCMHandle, m_HWParams, SND_PCM_FORMAT_S16_LE) < 0)
  {
    std::cout << "[ALSA] Cannot set sample format, " << stderr << std::endl;
  }

  if (snd_pcm_hw_params_set_rate_near(m_PCMHandle, m_HWParams, (unsigned int*)&audioFile.sample_rate, 0) < 0)
  {
    std::cout << "[ALSA] Cannot set sample rate, " << stderr << std::endl;
  }

  if (snd_pcm_hw_params_set_channels(m_PCMHandle, m_HWParams, audioFile.channels) < 0)
  {
    std::cout << "[ALSA] Cannot set channel count, " << stderr << std::endl;
  }

  if (snd_pcm_hw_params(m_PCMHandle, m_HWParams) < 0)
  {
    std::cout << "[ALSA] Cannot set parameters, " << stderr << std::endl;
  }

  if (snd_pcm_prepare(m_PCMHandle) < 0)
  {
    std::cout << "[ALSA] Cannot prepare audio interface for use, " << stderr << std::endl;
  }
}

void AudioSink::PlayAudioData(const unsigned char* buffer, AudioFile audioFile, size_t bytesDecoded)
{
  int frames = bytesDecoded / (audioFile.channels * 2); // 2 bytes per sample since we use S16 -> 16bits = 2bytes per channel
  int err = snd_pcm_writei(m_PCMHandle, (short*)buffer, frames);

  std::cout << "Bytes: " << bytesDecoded
            << " Frames: " << frames << std::endl;
  if (err == -EPIPE)
  {
    // this if statement checks whether the audio playback breaked because data wasn't delivered at time;
    // if so, continue playing the audio (prepare/reset the PCM device and try to play audio again)
    snd_pcm_prepare(m_PCMHandle);
  }

  if (err < 0)
  {
      snd_pcm_prepare(m_PCMHandle);
  }
}

void AudioSink::Shutdown()
{
  snd_pcm_drain(m_PCMHandle);
  snd_pcm_close(m_PCMHandle);
}
