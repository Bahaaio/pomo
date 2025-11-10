// Package ascii provides ASCII art style.
package ascii

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	zero = `
   ▄▄▄▄   
  ██▀▀██  
 ██    ██ 
 ██ ██ ██ 
 ██    ██ 
  ██▄▄██  
   ▀▀▀▀   
`
	one = `
   ▄▄▄    
  █▀██    
    ██    
    ██    
    ██    
 ▄▄▄██▄▄▄ 
 ▀▀▀▀▀▀▀▀ 
`
	two = `
  ▄▄▄▄▄   
 █▀▀▀▀██▄ 
       ██ 
     ▄█▀  
   ▄█▀    
 ▄██▄▄▄▄▄ 
 ▀▀▀▀▀▀▀▀ 
`
	three = `
  ▄▄▄▄▄   
 █▀▀▀▀██▄ 
      ▄██ 
   █████  
      ▀██ 
 █▄▄▄▄██▀ 
  ▀▀▀▀▀   
`
	four = `
     ▄▄▄  
    ▄███  
   █▀ ██  
 ▄█▀  ██  
 ████████ 
      ██  
      ▀▀  
`
	five = `
 ▄▄▄▄▄▄▄  
 ██▀▀▀▀▀  
 ██▄▄▄▄   
 █▀▀▀▀██▄ 
       ██ 
 █▄▄▄▄██▀ 
  ▀▀▀▀▀   
`
	six = `
   ▄▄▄▄   
  ██▀▀▀█  
 ██ ▄▄▄   
 ███▀▀██▄ 
 ██    ██ 
 ▀██▄▄██▀ 
   ▀▀▀▀   
`
	seven = `
 ▄▄▄▄▄▄▄▄ 
 ▀▀▀▀▀███ 
     ▄██  
     ██   
    ██    
   ██     
  ▀▀      
`
	eight = `
   ▄▄▄▄   
 ▄██▀▀██▄ 
 ██▄  ▄██ 
  ██████  
 ██▀  ▀██ 
 ▀██▄▄██▀ 
   ▀▀▀▀   
`
	nine = `
   ▄▄▄▄   
 ▄██▀▀██▄ 
 ██    ██ 
 ▀██▄▄███ 
   ▀▀▀ ██ 
  █▄▄▄██  
   ▀▀▀▀   
`
	colon = `

    ▄▄    
    ██    
          
    ██    
    ▀▀    

`
)

func ToASCIIArt(number string) string {
	digits := make([]string, 0, len(number))

	for _, digit := range number {
		digits = append(digits, digitToASCIIArt(digit))
	}

	asciiDigits := lipgloss.JoinHorizontal(lipgloss.Top, digits...)
	return asciiDigits
}

func digitToASCIIArt(digit rune) string {
	switch digit {
	case '0':
		return zero
	case '1':
		return one
	case '2':
		return two
	case '3':
		return three
	case '4':
		return four
	case '5':
		return five
	case '6':
		return six
	case '7':
		return seven
	case '8':
		return eight
	case '9':
		return nine
	case ':':
		return colon
	default:
		return ""
	}
}
