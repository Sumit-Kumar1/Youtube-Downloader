# Youtube Video Downloader

- This is a **HTMX-golang** project runs at port `http://localhost:12344/`
- This service requires ffmpeg library to be installed, please install ffmpeg first. Otherwise you will be getting lots of error for download every-time

Steps to run the project (makefile)

1. Clone the project
2. Make sure `make ` is installed in your system
3. Run the project using `make` help:
    3.1 to build simply `make build`
    3.2 to run `make run`
4. Open the browser and go to `http://localhost:12344/`
5. You can first get the info about the video/playlist/youtube-shorts by entering the youtube-url in the input box and clicking the "Get Info" button
6. You can also download the video by clicking the "Download" button in the page for which download is requested
7. You can track download progress by using `/getStatus` endpoint hitting a button in page for which download is requested

## Enhancements

- [x] Video Quality Selection:
  - Allow users to choose different video quality options (low, medium, high)
  - Download videos in the selected quality using the YouTube Data API

- [x] Audio Download:
  - Provide an option to download audio instead of video

- [ ] Audio Download, music info in downloaded audios e.g. Thumbnail, track info etc.

- [ ] Subtitles Download:
  - Download and display subtitles for videos
  - Allow users to download subtitles in various formats (e.g., SRT, VTT)

- [ ] Progress Tracking:
  - Display a progress bar to indicate the download progress
  - Update the progress bar in real-time using JavaScript

- [ ] Error Handling: (WIP)
  - Handle errors gracefully and provide informative messages to users
  - Log errors for debugging purposes

- [ ] User Interface Enhancements: (WIP)
  - Use HTML5 audio and video tags to embed the downloaded video and audio files
  - Add a search bar to allow users to search for videos

- [ ] Playlist Download:
  - Allow users to download entire playlists
  - Use the YouTube Data API to retrieve the list of videos in the playlist

- [ ] Database Storage:
  - Store downloaded videos and metadata in a database
  - This can allow users to manage their downloaded videos and search for them later

### Additional Features

- [ ] Support for multiple video formats (e.g., MP4, AVI)
- [ ] Video conversion to different formats
- [ ] Batch video download
- [ ] Automatic video thumbnail generation
- [ ] Social media sharing options
