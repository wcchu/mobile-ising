suppressPackageStartupMessages(library(ggplot2))

state_hist <- read.csv('../model/state_hist.csv')

mag_hist <- read.csv('../model/mag_hist.csv')

mag_plot <-
  ggplot(mag_hist) +
  geom_point(aes(x = time, y = mag))
print(mag_plot)
