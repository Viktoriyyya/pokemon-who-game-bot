package pokemonwhogame

import (
	"fmt"
	"image"
	"image/color"
	"net/http"
)

//Pokemon holds only the necessary info for the game to work
type Pokemon struct {
	id   int
	name string
	img  image.Image
}

//PokemonList is just the pokemon CSV
type PokemonList [][]string

const pokemonAssets = "https://assets.pokemon.com/assets/cms2/img/pokedex/full/%03d.png"

func (p PokemonList) getPokemon(id int) Pokemon {
	resp, _ := http.Get(fmt.Sprintf(pokemonAssets, id))
	decodedImage, _, _ := image.Decode(resp.Body)
	return Pokemon{id, p[id][30], decodedImage}
}

//shadowImage is the "shadow" version of an image: all non-alpha pixels are changed to black.
type shadowImage struct {
	originalImage image.Image
}

func (i shadowImage) ColorModel() color.Model {
	return i.originalImage.ColorModel()
}

func (i shadowImage) Bounds() image.Rectangle {
	return i.originalImage.Bounds()
}

func (i shadowImage) At(x, y int) color.Color {
	_, _, _, a := i.originalImage.At(x, y).RGBA()
	if a == 0 {
		return color.RGBA{255, 255, 255, 255}
	}
	return color.RGBA{0, 0, 0, 255}
}
