FROM scratch

COPY ./backend/cmd/laga/laga /laga

ENTRYPOINT [ "/laga" ]
