use std::fs::File;
use std::io::BufReader;
use rodio::{Decoder, MixerDeviceSink, source::Source};


fn main() {
    let sink_handle = rodio::DeviceSinkBuilder::open_default_sink()
        .expect("open default audio stream");
    let player = rodio::Player::connect_new(&sink_handle.mixer());

    let file = BufReader::new(File::open("../data/music/yippee-tbh.mp3").unwrap());
    let source = rodio::Decoder::try_from(file).unwrap();

    player.append(source);
    player.sleep_until_end();
}
