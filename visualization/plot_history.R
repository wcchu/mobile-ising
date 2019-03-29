suppressPackageStartupMessages(library(tidyverse))

## state history
state_hist <- read.csv('../model/state_hist.csv')
state_plot <-
  ggplot(state_hist) +
  geom_point(aes(x = x, y = y, color = spin, size = conns)) +
  facet_grid(temp ~ time)
ggsave("state_hist.png", state_plot, width = 100, height = 10, units = "cm")

## magnetization history
macro_hist <-
  read.csv('../model/macro_hist.csv') %>%
  gather(key = "key", value = "value", c(mag, ener))

macro_plot <-
  ggplot(macro_hist) +
  geom_point(aes(x = time, y = value), size = 0.2) +
  facet_grid(temp ~ key)
ggsave("macro_hist.png", macro_plot, width = 15, height = 50, units = "cm")

