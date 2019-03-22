suppressPackageStartupMessages(library(tidyverse))

## state history
state_hist <- read.csv('../model/state_hist.csv')
state_plot <-
  ggplot(state_hist) +
  geom_point(aes(x = x, y = y, color = spin)) +
  facet_grid(time ~ temp)
print(state_plot)

## magnetization history
mag_hist <- read.csv('../model/mag_hist.csv')
mag_plot <-
  ggplot(mag_hist) +
  geom_point(aes(x = time, y = mag)) +
  facet_grid(. ~ temp)
print(mag_plot)
