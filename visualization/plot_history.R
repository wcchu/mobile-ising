suppressPackageStartupMessages(library(tidyverse))

## state history
state_hist <- read.csv('../model/state_hist.csv')
ntemps <- length(unique(state_hist$temp))
state_plot <-
  ggplot(state_hist) +
  geom_point(aes(x = x, y = y, color = as.character(spin)), size = 0.1) +
  facet_grid(temp ~ time)
ggsave("state_hist.png", state_plot, width = 25, height = 2*ntemps, units = "cm")

## magnetization history
macro_hist <-
  read.csv('../model/macro_hist.csv') %>%
  gather(key = "key", value = "value", c(mag, ener))

macro_plot <-
  ggplot(macro_hist) +
  geom_point(aes(x = time, y = value), size = 0.2) +
  facet_grid(temp ~ key)
ggsave("macro_hist.png", macro_plot, width = 15, height = 5*ntemps, units = "cm")

