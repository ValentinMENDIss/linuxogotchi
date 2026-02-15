#include "AudioDecoder.h"
#include <iostream>
#include <mpg123.h>


void Mpg123Decoder::OpenAudioFile(std::string filepath)
{
  if (mpg123_open(m_Mpg123Handle, filepath.c_str()) || mpg123_getformat(m_Mpg123Handle, &m_AudioFile.sample_rate, &m_AudioFile.channels, &m_AudioFile.encoding) != MPG123_OK)
    std::cout << "[MPG123] Error opening MPG123 decoder" << mpg123_strerror(m_Mpg123Handle) << std::endl;

  // force mpg123 to use 16-bit signed if needed (since ALSA uses S16 format)
  mpg123_format_none(m_Mpg123Handle);
  mpg123_format(m_Mpg123Handle, m_AudioFile.sample_rate, m_AudioFile.channels, MPG123_ENC_SIGNED_16);
}

bool Mpg123Decoder::ReadAudioFile()
{
  int result = mpg123_read(m_Mpg123Handle, m_Buffer, sizeof(m_Buffer), &m_BytesDecoded);
  if (result == MPG123_DONE)
    return false;

  if (result != MPG123_OK)
  {
    std::cout << "[MPG123] Error while trying to read audio file, " << mpg123_strerror(m_Mpg123Handle) << std::endl;
    return false;
  }

  return true;
}

void Mpg123Decoder::Shutdown()
{
  mpg123_close(m_Mpg123Handle);
  mpg123_delete(m_Mpg123Handle);
  mpg123_exit();
}
