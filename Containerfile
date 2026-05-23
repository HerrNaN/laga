FROM scratch

COPY ./backend/laga /laga

ENTRYPOINT [ "/laga" ]
