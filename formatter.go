package logs

type Formatter interface  {
	Format(*DataPkg) ([]byte, error)
}

const (
	color_default 	= "\033[0m"
	color_red 	= "\033[31m"
	color_green 	= "\033[32m"
	color_yello	= "\033[33m"
	color_blue 	= "\033[34m"
	color_attr	= "\033[36m"


)
