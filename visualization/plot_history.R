suppressPackageStartupMessages(library(tidyverse))

## state history
state_hist <- read.csv('../model/state_hist.csv')
state_plot <-
  ggplot(state_hist) +
  geom_point(aes(x = x, y = y, color = spin)) +
  facet_grid(temp ~ time)
print(state_plot)

## magnetization history
macro_hist <-
  read.csv('../model/macro_hist.csv') %>%
  gather(key = "key", value = "value", c(mag, ener))

macro_plot <-
  ggplot(macro_hist) +
  geom_point(aes(x = time, y = value), size = 0.5) +
  facet_grid(key ~ temp)
print(macro_plot)
