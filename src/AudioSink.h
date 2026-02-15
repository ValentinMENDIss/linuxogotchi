#pragma once

#include <alsa/asoundlib.h>
#include "AudioDecoder.h"

#include <iostream>

class AudioSink
{
private:
  snd_pcm_t* m_PCMHandle;
  snd_pcm_hw_params_t* m_HWParams;
public:
  void PrepareAudio(AudioFile audioFile);
  void PlayAudioData(const unsigned char* buffer, AudioFile audioFile, size_t bytesDecoded);
  void Shutdown();
};
