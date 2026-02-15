#include "MusicPlayer.h"
#include "AudioDecoder.h"
#include "AudioSink.h"
#include <memory>

MusicPlayer::MusicPlayer()
{
  m_AudioDecoder = std::make_unique<Mpg123Decoder>();
  m_AudioSink = std::make_unique<AudioSink>();
}

void MusicPlayer::PlayMusic(std::string filepath)
{
  m_AudioDecoder->OpenAudioFile(filepath);
  m_AudioSink->PrepareAudio(m_AudioDecoder->GetAudioData());

  while (m_AudioDecoder->ReadAudioFile())
  {
    size_t bytesDecoded = m_AudioDecoder->GetBytesDecoded();
    if (bytesDecoded == 0) {
        std::cout << "Warning: Read succeeded but 0 bytes decoded. Breaking to avoid potential infinite loop." << std::endl;
        break;
    }

    m_AudioSink->PlayAudioData(m_AudioDecoder->GetBuffer(), m_AudioDecoder->GetAudioData(), bytesDecoded);
  }
  m_AudioDecoder->Shutdown();
  m_AudioSink->Shutdown();
}
