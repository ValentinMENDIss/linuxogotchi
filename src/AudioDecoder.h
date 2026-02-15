#pragma once

#include <alsa/asoundlib.h>
#include <mpg123.h>

#include <iostream>

struct AudioFile
{
  long sample_rate;
  int channels;
  int encoding;
};

class AudioDecoder
{
  private:
  public:
    virtual void OpenAudioFile(std::string path) {}
    virtual bool ReadAudioFile() { return false; }
};

class Mpg123Decoder : AudioDecoder
{
  private:
  int m_Err = MPG123_OK;

  AudioFile m_AudioFile;
  mpg123_handle* m_Mpg123Handle = mpg123_new(nullptr, &m_Err);

  unsigned char m_Buffer[8192];
  size_t m_BytesDecoded;
  public:

    void OpenAudioFile(std::string filepath) override;
    bool ReadAudioFile() override;
    void Shutdown();

    inline const unsigned char* GetBuffer() { return m_Buffer; }
    inline const AudioFile GetAudioData()   { return m_AudioFile; }
    inline const size_t GetBytesDecoded()    { return m_BytesDecoded; }
};
