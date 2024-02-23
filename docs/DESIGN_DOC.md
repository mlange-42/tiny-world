# Tiny World - Design Document

## Game Idea

A tiny colony of city builder.

Players grow their own small 2D iso landscape by adding terrain tiles.
On the terrain tiles, "land use" can be placed.

A simple economy and demography are simulated.

## Game Elements

### Terrain

* Air (nothing)
* Grass/land
* Desert
* Water

### Land Use / Buildings

Natural land use

* Air (nothing)
* Trees
* Rocks

Other land use

* Path
* Fields

Buildings

* Farm
* Fisherman
* Lumberjack
* Mason

Potentially, later:

* Sawmill
* Wind mill
* Water mill
* Bakery
* Residential

### Production

Production buildings need an adjacent road in one of the 4 neighboring tiles.
They produce one unit per minute for each relevant land use tile in the 8 neighboring tiles.

Further, food storage must be positive for any production except food

* field -> farm -> food
* water -> fisherman -> food
* trees -> lumberjack -> wood
* rocks -> mason -> stones

### Consumption

Production buildings require 5 units of food per minute.
Food production buildings require only 1 unit per minute.
