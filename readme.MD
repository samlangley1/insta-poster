# Project Title

insta-post

## Description

Posts random images from a local directory to Instagram. Works on Windows & Linux systems


## Requirements

- Docker
- An Instagram Account

## Usage

The application should always be run as a docker container and the following environment variables are required. The easiest way is by providing them in a .env file in the root of the project.

```
- touch .env
```

| Name      | Description | Required |
| ----------- | ----------- | ----------- |
| INSTAGRAM_USERNAME    |   Username required to log into Instagram   | Yes |
| INSTAGRAM_PASSWORD  |   Password required to log into Instagram   | Yes |
| POST_CAPTION  |    Caption used when posting content to Instagram     | No |
| CAPTION_SUMMARY | *Advanced* Used for generating random captions from a set of provided values [*See docs*](#randomising-captions) | No |
| CAPTION_CTA | *Advanced* Used for generating random captions from a set of provided values [*See docs*](#randomising-captions) | No |
| CAPTION_HASHTAG | *Advanced* Used for generating random captions from a set of provided values [*See docs*](#randomising-captions) | No |
| HTTP_PROXY | Proxy address if routing requests through a proxy | No |


The image content can be stored in any folder, as long as you have permission to bindmount it to the Docker container.

The image directory MUST have a subdirectory matching the Instagram account name exactly, for example, if your images at stored at ~/Documents/insta-poster/images then your structure should be: ~/Documents/insta-poster/images/*account-name*/image.jpg
You can store multiple accounts images in the same place, but each account must have it's own subdirectory.

```
- docker build -t insta-poster:latest .
- docker run -v image-directory-path:/opt/images insta-poster:latest
```

In the above example, your image directory path would be ~/Documents/insta-poster/images

For Windows users, make sure you use the \ syntax for filepath, e.g C:\Users\John\Documents\insta-poster\images


## Randomising Captions

If you're looking to mix up your captions every post, rather than providing a static caption, there is functionality to build a caption from 3 sets of values, a value from each of the sets is generated at random and then built as a final caption.

Each caption is split into 3 sections, take the following caption for example:

Double tap to show your appreciation\
\-----\
Follow us for more amazing content!\
\-----\
#2023 #f4f #amazing #fashion


Summary: This is the first line, used to provide context or summarise the post

CTA: The middle section, used to provide a simple call to action

Hashtags: Hashtags separated by a space

To generate a random caption, three additional environment variables are required, as well as leaving the "POST_CAPTION" variable blank or omitting it completely.

CAPTION_SUMMARY = Comma separated list of summaries. One is used for each post. Example:
Double tap to show your appreciation,Leave a like if you agree with this post!,Send this to 2 friends!

CAPTION_CTA = Comma separated list of CTA's. One is used for each post. Example:
Follow us for more amazing content!,Why not share with a friend?,For more daily content give us a follow!

CAPTION_HASHTAG = Comma separated list of hashtags. 15 are selected randomly for each post. Example:
#2023,#f4f,#amazing,#fashion,#l4l,#wow,#boy,#girl,#viral,#foryou,#fyp,#explore,#style,#instagood,#ootd,#outfit,#fitcheck,#today,#yesterday,#tomorrow
