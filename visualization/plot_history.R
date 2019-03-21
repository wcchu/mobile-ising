suppressPackageStartupMessages(library(tidyverse))

## state history
state_hist <- read.csv('../model/state_hist.csv')
state_plot <-
  ggplot(state_hist %>% filter(time == 0.0)) +
  geom_point(aes(x = x, y = y, color = spin))
print(state_plot)

## magnetization history
mag_hist <- read.csv('../model/mag_hist.csv')
mag_plot <-
  ggplot(mag_hist) +
  geom_point(aes(x = time, y = mag))
print(mag_plot)
