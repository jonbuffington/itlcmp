include $(GOROOT)/src/Make.inc

TARG=itlcmp
GOFILES=\
  library.go\
	main.go\
  parser.go\
  statistics.go\

include $(GOROOT)/src/Make.cmd
