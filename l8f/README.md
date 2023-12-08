# lyrics818 format (l8f)
a lyrics video format written in golang 

## Features
* Two video tracks: one for laptop size and the other for mobile size
* Written in Golang

## Description of the Format

The general description looks like this:
```
{header_length}
{header}
{audio}
{laptop_frames_lump}
{mobile_frames_lump}
```

This format uses a framerate of 1

### Description of {header} section

The `{header}` section made up of some subsections. It looks like this

```
laptop_unique_frames:
{number}: {size}
{number}: {size}
::
laptop_frames:
{frame_number}: {unique_frame_number}
{frame_number}: {unique_frame_number}
::
mobile_unique_frames:
{number}: {size}
{number}: {size}
::
mobile_frames:
{frame_number}: {unique_frame_number}
{frame_number}: {unique_frame_number}
::
binary:
audio: {audio_size_bytes}
laptop_frames_lump: {laptop_video_size_bytes}
mobile_frames_lump: {mobile_video_size_bytes}
::
```


### Description of the {audio} section

The `{audio}` section takes mp3 data and writes it unparsed to the video format

### Description of the {laptop_frames_lump} / {mobile_frames_lump} sections

The `{video}` section is made up of a lump file of unique frames.
The unique frames must be of **png** format
