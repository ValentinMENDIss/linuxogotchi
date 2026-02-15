#pragma once

#include "AudioDecoder.h"
#include "AudioSink.h"
#include <alsa/asoundlib.h>
#include <memory>

class MusicPlayer
{
  private:
    std::unique_ptr<Mpg123Decoder> m_AudioDecoder;
    std::unique_ptr<AudioSink> m_AudioSink;
  public:
    MusicPlayer();
    void PlayMusic(std::string filepath);
};
