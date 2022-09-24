# srtmod

**srtmod** is a CLI application written in [Go](https://go.dev/) and delivered as a [single executable](https://github.com/markhaur/srtmod/releases/tag/v0.1.1) to add/subtract time offset from .srt subtitle files.
## Usage
Here's a sample `config.yml` file.
```yml
files:
  - inputFile: "input.srt"
    outputFile: "output.srt"
    offset: -1m
  - inputFile: "input.srt"
    outputFile: "output.srt"
    offset: 5s
  - inputFile: "input.srt"
    outputFile: "output.srt"
    offset: -1m4s
```
You can process batch of .srt files through `srtmod` by using `config.yml` file
```bash
./srtmod -config=config.yml
```
or can process single .srt file using args
```bash
./srtmod -i=input.srt -o=output.srt -offset=-1m5s
```
## Contributing
Want to contribute? Awesome! The most basic way to show your support is to star the project or to raise issues. Pull requests are highly welcome.

Please make sure to update tests as appropriate.
## Potential Maintainers
[Jarri Abidi](https://github.com/jarri-abidi)\
[Faisal Nisar](https://github.com/markhaur)
## Credits
This project exists, thanks to all the people who contribute.

<a href="https://github.com/codeforcauseorg/edu-client/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=markhaur/srtmod" />
</a>
