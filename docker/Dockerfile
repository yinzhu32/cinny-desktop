FROM archlinux:latest

RUN pacman -Sy && pacman -S --noconfirm --noprogressbar gtk4 webkit2gtk gobject-introspection go gcc make && mkdir /source

WORKDIR /source
