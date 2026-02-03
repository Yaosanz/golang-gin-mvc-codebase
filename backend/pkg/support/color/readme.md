## Introduction
The ```color``` package provides a set of functions to colorize the output of the terminal using **PTerm** library.

## Specific Color
The package provides methods to create printers for specific colors. These methods allow you to easily colorize terminal output.

- `color.Red()`
- `color.Green()`
- `color.Yellow()`
- `color.Blue()`
- `color.Magenta()`
- `color.Cyan()`
- `color.White()`
- `color.Black()`
- `color.Gray()`
- `color.Default()`

```
import "<app-name>/pkg/support/color"

color.Blue().Println("Hello, World!")
color.Green().Printf("Hello, %s!", "World")
```

## Printer Methods
A contracts/support.Printer provides the following methods to print or format text with color:

- `Print` - Print text
- `Println` - Print text with a new line
- `Printf` - Print formatted text
- `Sprint` - Return colored text
- `Sprintln` - Return colored text with a new line
- `Sprintf` - Return formatted colored text

## Custom Color
`color.New`

The `color.New` function creates a new color printer. You can use this object to colorize the output of the terminal.

```
import "<app-name>/pkg/support/color"

color.New(color.FgRed).Println("Hello, World!")
```