# My CV Project

## Table of Contents

- [Overview](#overview)
- [Dependencies](#dependencies)
- [Using Docker as an alternative](#using-docker-as-an-alternative)
- [Generating the PDF Document](#generating-the-pdf-document)
- [Inspirations](#inspirations)

## Overview

This project parses my [CV](./data/cv.json) that is written as
a JSON document and generates a PDF document.
This project contains a Go application and uses the ConTeXt document processor
which are used to parse and generate a PDF document from the CV.
I chose to use ConTeXt as the document processor because it gives me more
control to edit page layouts, configure fonts,
perform additional formatting and add custom functionality.

Support for other formats such as HTML will be available soon.

### View/Download the CV

The latest build of the generated PDF document can be downloaded [here](__output/cv.pdf)

## Dependencies

If you are interested in generating the PDF document on your own machines
or want to use the project to generate your own CV, then below is a list
of dependencies that you'll need to install:

- **Go** - Please go [here](https://golang.org/dl/) to download the latest version of the Go programming language.
- **ConTeXt** - I recommend installing version 1.02 or higher. Please go [here](https://wiki.contextgarden.net/ConTeXt_Standalone) for installation instructions for your distribution.
- **The Carlito font (ttf-carlito)** - In prevous iterations of my CV I used the Calibri font. Carlito is the free, metric compatible alternative to this and is specified in the TEX template.
  - For Ubuntu/Debian installation you can use `apt`:
    ```bash
    $ apt-get install font-crosextra-carlito
    ```
  - For Arch Linux you can use `pacman`:
    ```bash
    $ pacman -S ttf-carlito
    ```
  - Alternatively you can download the font from https://fontlibrary.org/en/font/carlito
  - Once this font is installed you'll need to update ConTeXt so it can find the font when generating the PDF:
    ```bash
    $ OSFONTDIR=/usr/share/fonts
    $ mtxrun --script fonts --reload
    ```
- **Make** - This should be available on all Unix distributions.

Once the dependencies are installed you can follow the
[Generating the PDF Document](#generating-the-pdf-document) section below.

## Using Docker as an alternative

If you prefer not to install the dependencies above,
I have created a Docker image installed with the above dependencies.
This image is used to build and publish the CV via the GitLab CI pipleines.
The image is built using this [Dockerfile](./docker/Dockerfile) and is
pullable from [GitLab's container registry](https://gitlab.com/dananglin/cv/container_registry).

To use the image follow the steps below:

1. Make sure you've clone this project to your workspace:
    ```bash
    $ git clone https://gitlab.com/dananglin/cv.git
    $ cd cv
    ```

2. Create the docker container and mount the current directory to the /project directory inside the container:
    ```bash
    $ docker run --rm -it -v ${PWD}:/project registry.gitlab.com/dananglin/cv/cv-builder:master-5fbdaa5a bash

    # Once inside the docker container
    $ cd /project
    ```

3. Follow the [Generating the PDF Document](#generating-the-pdf-document) section below.

## Generating the PDF Document

The PDF document can be generated by running the following command:

```bash
$ make pdf
```

The PDF generation is completed in two steps:

1. The Go application will generate a TEX file using the JSON document and the TEX template files located in the [template](./template) directory.
2. ConTeXt is then used to generate the PDF document from the generated TEX file.

## Inspirations

- [The Markdown Resume](https://mszep.github.io/pandoc_resume/) - This project uses ConTeXt and pandoc to convert Markdown based CVs into multiple formats including PDF, HTML and DOCX. This is where I discovered ConTeXt.
- [melkir/resume](https://github.com/melkir/resume) - This project generates CVs using Go and LaTeX.
