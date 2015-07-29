Simple tex calendar generator
======

Generates a simple monday-sunday (full weeks) calendar using a tex template.

Preview: [English](preview_en.pdf) [German](preview_de.pdf)

Get the generator
-----
```bash
$ go get github.com/blang/texcal
```

Get one of the templates (calendar*.tex) and save it as `calendar.tex`.

Usage
-----
```
./texcal -output ./output.tex -begin="06.03.2015" -days 30
./texcal -output ./output.tex -begin="06.03.2015" -end="15.04.2015"
./texcal -begin="06.03.2015" > 30days.tex
./texcal > 30days.tex
```

Other languages
------

```
./texcal -lang="de" -tmpl="calendar_german.tex" > german.tex
```


To PDF
------
```
pdflatex  -shell-escape ./output.tex
```

License
-----

See [LICENSE](LICENSE) file.