package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"image/color"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"strconv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"time"
)

var gameStarted int 
var mainGameController gameController
var normalFont font.Face
const dpi = 60

//Struct of the semaphores
type semaphore struct {
	//Stop = 1
	//Go = 0
	color int
	//Timer to change color
	counter int
	//Positio of the semaphore
	positionX float64
	positionY float64
	image *ebiten.Image
	options *ebiten.DrawImageOptions
}

//Struct of the cars
type car struct {
	speedX       float64
	speedY       float64
	startingSpeedX float64
	startingSpeedY float64
	carRoute	 route
	//0 car is on route
	//1 car arrived
	status int
	positionX float64
	positionY float64
	routeNumber int
	rotationType int
	image *ebiten.Image
	options ebiten.DrawImageOptions
	lastStep bool
	alreadyStopped bool
	negative bool
	//Time taken to car to start route
	startTime int
}

type route struct {
	startX float64
	startY float64
	destinationX float64
	destinationY float64
}

type traficLightController struct {
	lightsImage        *ebiten.Image
	lightsOptions      *ebiten.DrawImageOptions
	positionX int
	positionY int
	rotationType int
}

//Struct of the game controller
type gameController struct {
	cars               []car
	semaphores         []semaphore
	traficLightControllers []traficLightController
	numberOfCars       int
	numberOfSemaphores int
	screenWidth int 
	screenHeight int
	startingPositionsX []float64
	startingPositionsY []float64
	endingPositionsX []float64
	endingPositionsY []float64
	stopsPositionX []float64
	stopsPositionY []float64
	//0 = X
	//1 = Y
	stopsType []int
}

//Function of semaphore
func semaphoreBehavior(carIndex int) {
	for {
		//Changin of color
		if mainGameController.semaphores[carIndex].counter == 6 {
			if mainGameController.semaphores[carIndex].color == 0 {
				mainGameController.semaphores[carIndex].color = 1
			} else if mainGameController.semaphores[carIndex].color == 1 {
				mainGameController.semaphores[carIndex].color = 0
			}
			mainGameController.semaphores[carIndex].counter = 0
		} else {
			mainGameController.semaphores[carIndex].counter++
			time.Sleep(1 * time.Second)
		}
	}
}

//Function of car
func carBehavior(carIndex int) {
	for {
		if(gameStarted == 3){
			for j := mainGameController.cars[carIndex].startTime; j > 0; j-- {
				mainGameController.cars[carIndex].startTime --;
				time.Sleep(1 * time.Second)
			}
	
			if(mainGameController.cars[carIndex].alreadyStopped == true){
				//Slow car moving in x axis
				if(mainGameController.cars[carIndex].speedX != 0 && mainGameController.cars[carIndex].speedY == 0){
					tempOldVelocity := mainGameController.cars[carIndex].speedX
					if(mainGameController.cars[carIndex].speedX>0){
						mainGameController.cars[carIndex].speedX = mainGameController.cars[carIndex].speedX - (mainGameController.cars[carIndex].speedX * .90)
					}else{
						mainGameController.cars[carIndex].speedX = mainGameController.cars[carIndex].speedX + (math.Abs(mainGameController.cars[carIndex].speedX) * .90)
					}
					
					time.Sleep(3 * time.Second)
					mainGameController.cars[carIndex].speedX = tempOldVelocity
				}else if(mainGameController.cars[carIndex].speedX == 0 && mainGameController.cars[carIndex].speedY != 0){
					tempOldVelocity := mainGameController.cars[carIndex].speedY
					if(mainGameController.cars[carIndex].speedY>0){
						mainGameController.cars[carIndex].speedY = mainGameController.cars[carIndex].speedY - (mainGameController.cars[carIndex].speedY * .95)
					}else{
						mainGameController.cars[carIndex].speedY = mainGameController.cars[carIndex].speedY + (math.Abs(mainGameController.cars[carIndex].speedY) * .95)
					}
					time.Sleep(3 * time.Second)
					mainGameController.cars[carIndex].speedY = tempOldVelocity
				}
			}

		}
	}
}

func drawBoard(screen *ebiten.Image){
	//Every sprite is 80*80 for moving the road just add or remove 80
	roadMain, _, _ := ebitenutil.NewImageFromFile("roads/roadNEWS.png", ebiten.FilterDefault)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(mainGameController.screenWidth/2)-40, float64(mainGameController.screenHeight/2)-40)
	screen.DrawImage(roadMain,op)

	road1, _, _ := ebitenutil.NewImageFromFile("roads/roadEW.png", ebiten.FilterDefault)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(float64(mainGameController.screenWidth/2)+40, float64(mainGameController.screenHeight/2)-40)
	screen.DrawImage(road1,op2)

	road2, _, _ := ebitenutil.NewImageFromFile("roads/roadEW.png", ebiten.FilterDefault)
	op3 := &ebiten.DrawImageOptions{}
	op3.GeoM.Translate(float64(mainGameController.screenWidth/2)-120, float64(mainGameController.screenHeight/2)-40)
	screen.DrawImage(road2,op3)

	road3, _, _ := ebitenutil.NewImageFromFile("roads/roadNS.png", ebiten.FilterDefault)
	op4 := &ebiten.DrawImageOptions{}
	op4.GeoM.Translate(float64(mainGameController.screenWidth/2)-40, float64(mainGameController.screenHeight/2)+40)
	screen.DrawImage(road3,op4)

	road4, _, _ := ebitenutil.NewImageFromFile("roads/roadNS.png", ebiten.FilterDefault)
	op5 := &ebiten.DrawImageOptions{}
	op5.GeoM.Translate(float64(mainGameController.screenWidth/2)-40, float64(mainGameController.screenHeight/2)-120)
	screen.DrawImage(road4,op5)

	road5, _, _ := ebitenutil.NewImageFromFile("roads/roadSE.png", ebiten.FilterDefault)
	op6 := &ebiten.DrawImageOptions{}
	op6.GeoM.Translate(float64(mainGameController.screenWidth/2)-40, float64(mainGameController.screenHeight/2)-200)
	screen.DrawImage(road5,op6)

	road6, _, _ := ebitenutil.NewImageFromFile("roads/roadSE.png", ebiten.FilterDefault)
	op7 := &ebiten.DrawImageOptions{}
	op7.GeoM.Translate(float64(mainGameController.screenWidth/2)-200, float64(mainGameController.screenHeight/2)-40)
	screen.DrawImage(road6,op7)

	road8, _, _ := ebitenutil.NewImageFromFile("roads/roadSW.png", ebiten.FilterDefault)
	op9 := &ebiten.DrawImageOptions{}
	op9.GeoM.Translate(float64(mainGameController.screenWidth/2)+120, float64(mainGameController.screenHeight/2)-40)
	screen.DrawImage(road8,op9)

	road10, _, _ := ebitenutil.NewImageFromFile("roads/roadNS.png", ebiten.FilterDefault)
	op11 := &ebiten.DrawImageOptions{}
	op11.GeoM.Translate(float64(mainGameController.screenWidth/2)-40, float64(mainGameController.screenHeight/2)+120)
	screen.DrawImage(road10,op11)

	road11, _, _ := ebitenutil.NewImageFromFile("roads/roadNW.png", ebiten.FilterDefault)
	op12 := &ebiten.DrawImageOptions{}
	op12.GeoM.Translate(float64(mainGameController.screenWidth/2)+200, float64(mainGameController.screenHeight/2)-200)
	screen.DrawImage(road11,op12)

	road13, _, _ := ebitenutil.NewImageFromFile("roads/roadEW.png", ebiten.FilterDefault)
	op14 := &ebiten.DrawImageOptions{}
	op14.GeoM.Translate(float64(mainGameController.screenWidth/2)+120, float64(mainGameController.screenHeight/2)-200)
	screen.DrawImage(road13,op14)

	road14, _, _ := ebitenutil.NewImageFromFile("roads/roadEW.png", ebiten.FilterDefault)
	op15 := &ebiten.DrawImageOptions{}
	op15.GeoM.Translate(float64(mainGameController.screenWidth/2)+40, float64(mainGameController.screenHeight/2)-200)
	screen.DrawImage(road14,op15)

	road15, _, _ := ebitenutil.NewImageFromFile("roads/roadNS.png", ebiten.FilterDefault)
	op16 := &ebiten.DrawImageOptions{}
	op16.GeoM.Translate(float64(mainGameController.screenWidth/2)+120, float64(mainGameController.screenHeight/2)+40)
	screen.DrawImage(road15,op16)

	road16, _, _ := ebitenutil.NewImageFromFile("roads/roadNE.png", ebiten.FilterDefault)
	op17 := &ebiten.DrawImageOptions{}
	op17.GeoM.Translate(float64(mainGameController.screenWidth/2)+120, float64(mainGameController.screenHeight/2)+120)
	screen.DrawImage(road16,op17)

	road17, _, _ := ebitenutil.NewImageFromFile("roads/roadEW.png", ebiten.FilterDefault)
	op18 := &ebiten.DrawImageOptions{}
	op18.GeoM.Translate(float64(mainGameController.screenWidth/2)+200, float64(mainGameController.screenHeight/2)+200)
	screen.DrawImage(road17,op18)

	road19, _, _ := ebitenutil.NewImageFromFile("roads/roadEW.png", ebiten.FilterDefault)
	op20 := &ebiten.DrawImageOptions{}
	op20.GeoM.Translate(float64(mainGameController.screenWidth/2)+200, float64(mainGameController.screenHeight/2)+120)
	screen.DrawImage(road19,op20)

	road20, _, _ := ebitenutil.NewImageFromFile("roads/roadNW.png", ebiten.FilterDefault)
	op21 := &ebiten.DrawImageOptions{}
	op21.GeoM.Translate(float64(mainGameController.screenWidth/2)-200, float64(mainGameController.screenHeight/2)+40)
	screen.DrawImage(road20,op21)

	road21, _, _ := ebitenutil.NewImageFromFile("roads/roadEW.png", ebiten.FilterDefault)
	op22 := &ebiten.DrawImageOptions{}
	op22.GeoM.Translate(float64(mainGameController.screenWidth/2)-280, float64(mainGameController.screenHeight/2)+40)
	screen.DrawImage(road21,op22)

	lights1, _, _ := ebitenutil.NewImageFromFile("roads/light.png", ebiten.FilterDefault)
	opl1 := &ebiten.DrawImageOptions{}
	opl1.GeoM.Translate(float64(mainGameController.screenWidth/2)-40, float64(mainGameController.screenHeight/2)-70)
	traficLightController1 := traficLightController{
		lightsImage: lights1,
		lightsOptions: opl1,
		positionX: (mainGameController.screenWidth/2)-40,
		positionY: (mainGameController.screenHeight/2)-70,
	}
	mainGameController.traficLightControllers = append(mainGameController.traficLightControllers, traficLightController1)
	mainGameController.semaphores[0].options = opl1
	mainGameController.semaphores[0].image = lights1

	lights2, _, _ := ebitenutil.NewImageFromFile("roads/light.png", ebiten.FilterDefault)
	opl2 := &ebiten.DrawImageOptions{}
	opl2.GeoM.Translate(float64(mainGameController.screenWidth/2)+125, float64(mainGameController.screenHeight/2)+40)
	traficLightController2 := traficLightController{
		lightsImage: lights2,
		lightsOptions: opl2,
		positionX: (mainGameController.screenWidth/2)+125,
		positionY: (mainGameController.screenHeight/2)+40,
	}
	mainGameController.traficLightControllers = append(mainGameController.traficLightControllers, traficLightController2)
	mainGameController.semaphores[1].options = opl2
	mainGameController.semaphores[1].image = lights2
	

	lights3, _, _ := ebitenutil.NewImageFromFile("roads/light.png", ebiten.FilterDefault)
	opl3 := &ebiten.DrawImageOptions{}
	opl3.GeoM.Translate(float64(mainGameController.screenWidth/2)-40, float64(mainGameController.screenHeight/2)+40)
	traficLightController3 := traficLightController{
		lightsImage: lights3,
		lightsOptions: opl3,
		positionX: (mainGameController.screenWidth/2)-40,
		positionY: (mainGameController.screenHeight/2)+40,
	}

	mainGameController.traficLightControllers = append(mainGameController.traficLightControllers, traficLightController3)
	mainGameController.semaphores[2].options = opl3
	mainGameController.semaphores[2].image = lights3

	lights4, _, _ := ebitenutil.NewImageFromFile("roads/light.png", ebiten.FilterDefault)
	opl4 := &ebiten.DrawImageOptions{}
	opl4.GeoM.Translate(float64(mainGameController.screenWidth/2)-200, float64(mainGameController.screenHeight/2)-20)
	traficLightController4 := traficLightController{
		lightsImage: lights4,
		lightsOptions: opl4,
		positionX: (mainGameController.screenWidth/2)-200,
		positionY: (mainGameController.screenHeight/2)-20,
	}
	mainGameController.traficLightControllers = append(mainGameController.traficLightControllers, traficLightController4)
	mainGameController.semaphores[3].options = opl4
	mainGameController.semaphores[3].image = lights4
	
	lights5, _, _ := ebitenutil.NewImageFromFile("roads/light2.png", ebiten.FilterDefault)
	opl5 := &ebiten.DrawImageOptions{}
	opl5.GeoM.Translate(float64(mainGameController.screenWidth/2)+120, float64(mainGameController.screenHeight/2)-170)
	traficLightController5 := traficLightController{
		lightsImage: lights5,
		lightsOptions: opl5,
		positionX: (mainGameController.screenWidth/2)+120,
		positionY: (mainGameController.screenHeight/2)-170,
	}
	mainGameController.traficLightControllers = append(mainGameController.traficLightControllers, traficLightController5)
	mainGameController.semaphores[4].options = opl5
	mainGameController.semaphores[4].image = lights5
	
	lights6, _, _ := ebitenutil.NewImageFromFile("roads/light2.png", ebiten.FilterDefault)
	opl6 := &ebiten.DrawImageOptions{}
	opl6.GeoM.Translate(float64(mainGameController.screenWidth/2)+40, float64(mainGameController.screenHeight/2)-10)
	traficLightController6 := traficLightController{
		lightsImage: lights6,
		lightsOptions: opl6,
		positionX: (mainGameController.screenWidth/2)+40,
		positionY: (mainGameController.screenHeight/2)-10,
	}
	mainGameController.traficLightControllers = append(mainGameController.traficLightControllers, traficLightController6)
	mainGameController.semaphores[5].options = opl6
	mainGameController.semaphores[5].image = lights6

	lights7, _, _ := ebitenutil.NewImageFromFile("roads/light2.png", ebiten.FilterDefault)
	opl7 := &ebiten.DrawImageOptions{}
	opl7.GeoM.Translate(float64(mainGameController.screenWidth/2)-50, float64(mainGameController.screenHeight/2)-10)
	traficLightController7 := traficLightController{
		lightsImage: lights7,
		lightsOptions: opl7,
		positionX: (mainGameController.screenWidth/2)-50,
		positionY: (mainGameController.screenHeight/2)-10,
	}
	mainGameController.traficLightControllers = append(mainGameController.traficLightControllers, traficLightController7)
	mainGameController.semaphores[6].options = opl7
	mainGameController.semaphores[6].image = lights7
	
	//Light1 -40 -40
	//Light2 200 40
	//Light3 -40 40
	//Light4 -200 -40
	//Light5 40 -200
	//Light6 -40 -40
	//Light6 -40 -40

	//Lights position and type
	mainGameController.stopsPositionX = append(mainGameController.stopsPositionX,float64(mainGameController.screenWidth/2)-40)
	mainGameController.stopsPositionY = append(mainGameController.stopsPositionY,float64(mainGameController.screenHeight/2)-70)
	mainGameController.stopsType = append(mainGameController.stopsType,1)

	mainGameController.stopsPositionX = append(mainGameController.stopsPositionX,float64(mainGameController.screenWidth/2)+125)
	mainGameController.stopsPositionY = append(mainGameController.stopsPositionY,float64(mainGameController.screenHeight/2)+40)
	mainGameController.stopsType = append(mainGameController.stopsType,1)

	mainGameController.stopsPositionX = append(mainGameController.stopsPositionX,float64(mainGameController.screenWidth/2)-40)
	mainGameController.stopsPositionY = append(mainGameController.stopsPositionY,float64(mainGameController.screenHeight/2)+40)
	mainGameController.stopsType = append(mainGameController.stopsType,1)

	mainGameController.stopsPositionX = append(mainGameController.stopsPositionX,float64(mainGameController.screenWidth/2)-200)
	mainGameController.stopsPositionY = append(mainGameController.stopsPositionY,float64(mainGameController.screenHeight/2)-20)
	mainGameController.stopsType = append(mainGameController.stopsType,1)
	 
	mainGameController.stopsPositionX = append(mainGameController.stopsPositionX,float64(mainGameController.screenWidth/2)+120)
	mainGameController.stopsPositionY = append(mainGameController.stopsPositionY,float64(mainGameController.screenHeight/2)-170)
	mainGameController.stopsType = append(mainGameController.stopsType,0)

	mainGameController.stopsPositionX = append(mainGameController.stopsPositionX,float64(mainGameController.screenWidth/2)+40)
	mainGameController.stopsPositionY = append(mainGameController.stopsPositionY,float64(mainGameController.screenHeight/2)-10)
	mainGameController.stopsType = append(mainGameController.stopsType,0)

	mainGameController.stopsPositionX = append(mainGameController.stopsPositionX,float64(mainGameController.screenWidth/2)-50)
	mainGameController.stopsPositionY = append(mainGameController.stopsPositionY,float64(mainGameController.screenHeight/2)-10)
	mainGameController.stopsType = append(mainGameController.stopsType,0)
}

//Game Loop
func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	//This part of the code draws the static road
	drawBoard(screen)

	//This part of the code draws all the cars in the screen and update the position of the cars
	for i := 0; i < mainGameController.numberOfCars; i++ {
		if(mainGameController.cars[i].startTime == 0){
			mainGameController.cars[i].options.GeoM.Translate(mainGameController.cars[i].speedX, mainGameController.cars[i].speedY)
			mainGameController.cars[i].positionX += mainGameController.cars[i].speedX
			mainGameController.cars[i].positionY += mainGameController.cars[i].speedY
			if mainGameController.cars[i].status == 1 {
				if(mainGameController.cars[i].speedY == 0.0){
					tempString := "Carro " + strconv.Itoa(i) + ": Ruta terminada  velocidad: " + fmt.Sprintf("%f", mainGameController.cars[i].speedX)
					text.Draw(screen, tempString , normalFont, 15 ,10*(i+1)*2, color.White)
				}else if(mainGameController.cars[i].speedX == 0.0){
					tempString := "Carro " + strconv.Itoa(i) + ": Ruta terminada  velocidad: " + fmt.Sprintf("%f", mainGameController.cars[i].speedY)
					text.Draw(screen, tempString , normalFont, 15 ,10*(i+1)*2, color.White)
				}
			}else if mainGameController.cars[i].status == 0 {
				if(mainGameController.cars[i].speedY == 0.0){
					tempString := "Carro " + strconv.Itoa(i) + ": En ruta  velocidad: " + fmt.Sprintf("%f", mainGameController.cars[i].speedX)
					text.Draw(screen, tempString , normalFont, 15 ,10*(i+1)*2, color.White)
				}else if(mainGameController.cars[i].speedX == 0.0){
					tempString := "Carro " + strconv.Itoa(i) + ": En ruta  velocidad: " + fmt.Sprintf("%f", mainGameController.cars[i].speedY)
					text.Draw(screen, tempString , normalFont, 15 ,10*(i+1)*2, color.White)
				}
			}

			//Route1
			//Decisions made in base of the position, rotation and destination of car
			if(mainGameController.cars[i].speedX > 0 && mainGameController.cars[i].speedY == 0.0 && mainGameController.cars[i].rotationType == 0 &&  mainGameController.cars[i].positionX >= 160.00 && mainGameController.cars[i].positionY <= 80.00 && mainGameController.cars[i].routeNumber == 0){
				mainGameController.cars[i].speedX = 0
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedY = -mainGameController.cars[i].startingSpeedX
				mainGameController.cars[i].rotationType = 1
			}
			if(mainGameController.cars[i].speedY <= 0 && mainGameController.cars[i].speedX == 0.0 && mainGameController.cars[i].rotationType == 1 && mainGameController.cars[i].positionX >= 160.00 && mainGameController.cars[i].positionY <= -80.00 && mainGameController.cars[i].routeNumber == 0){
				mainGameController.cars[i].speedX = mainGameController.cars[i].startingSpeedX
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedY = 0
				mainGameController.cars[i].rotationType = 0
			}
			if(mainGameController.cars[i].speedY == 0 && mainGameController.cars[i].speedX > 0.0 && mainGameController.cars[i].rotationType == 0 && mainGameController.cars[i].positionX >= 470.00 && mainGameController.cars[i].positionY <= -80.00 && mainGameController.cars[i].routeNumber == 0){
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*180))
				mainGameController.cars[i].rotationType = 3
				mainGameController.cars[i].speedX = 0
				mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedX
			}
			if(mainGameController.cars[i].speedY > 0 && mainGameController.cars[i].speedX == 0.0 && mainGameController.cars[i].rotationType == 3 && mainGameController.cars[i].positionX >= 470.00 && mainGameController.cars[i].positionY >= 80.00 && mainGameController.cars[i].routeNumber == 0){
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*180))
				mainGameController.cars[i].rotationType = 0
				mainGameController.cars[i].speedX = mainGameController.cars[i].startingSpeedX
				mainGameController.cars[i].speedY = 0
			}
			if(mainGameController.cars[i].speedY == 0 && mainGameController.cars[i].speedX > 0.0 && mainGameController.cars[i].rotationType == 0 && mainGameController.cars[i].positionX >= 580.00 && mainGameController.cars[i].positionY >= 80.00 && mainGameController.cars[i].routeNumber == 0){
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*180))
				mainGameController.cars[i].rotationType = 0
				mainGameController.cars[i].speedX = 0
				mainGameController.cars[i].speedY = 0
				mainGameController.cars[i].status = 1
			}
			//Route1

			//Route2
			//Decisions made in base of the position, rotation and destination of car
			if(mainGameController.cars[i].speedX == 0 && mainGameController.cars[i].speedY < 0.0 && mainGameController.cars[i].rotationType == 1 && mainGameController.cars[i].positionY <= -330.00 && mainGameController.cars[i].routeNumber == 1){
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = -mainGameController.cars[i].startingSpeedY
				mainGameController.cars[i].speedY = 0
				mainGameController.cars[i].rotationType = 0
			}
			
			if(mainGameController.cars[i].speedX > 0 && mainGameController.cars[i].speedY == 0.0 && mainGameController.cars[i].rotationType == 0 &&  mainGameController.cars[i].positionX >= 240.00 && mainGameController.cars[i].routeNumber == 1){
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = 0
				mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedY
				mainGameController.cars[i].rotationType = 3
			}

			if(mainGameController.cars[i].speedX == 0 && mainGameController.cars[i].speedY < 0.0 && mainGameController.cars[i].rotationType == 3 && mainGameController.cars[i].positionY <= -360.00 && mainGameController.cars[i].routeNumber == 1){
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = 0
				mainGameController.cars[i].speedY = 0
				mainGameController.cars[i].status = 1
			}
			//Route2

			//Route3
			//Decisions made in base of the position, rotation and destination of car
			if mainGameController.cars[i].lastStep == false && mainGameController.cars[i].speedX <= 0 && mainGameController.cars[i].speedY == 0.0 && mainGameController.cars[i].positionX <= -60 && mainGameController.cars[i].routeNumber == 2 {
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = 0
				mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedX
				mainGameController.cars[i].rotationType = 1
			}
			if mainGameController.cars[i].lastStep == false && mainGameController.cars[i].speedX == 0 && mainGameController.cars[i].speedY <= 0.0 && mainGameController.cars[i].rotationType == 1 && mainGameController.cars[i].positionY <= -160 && mainGameController.cars[i].routeNumber == 2 {
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = mainGameController.cars[i].startingSpeedX
				mainGameController.cars[i].speedY = 0
				mainGameController.cars[i].rotationType = 2
			}
			if mainGameController.cars[i].lastStep == false && mainGameController.cars[i].speedX <= 0 && mainGameController.cars[i].speedY == 0.0 && mainGameController.cars[i].rotationType == 2 && mainGameController.cars[i].positionX <= -235 && mainGameController.cars[i].routeNumber == 2 {
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = 0
				mainGameController.cars[i].speedY = -mainGameController.cars[i].startingSpeedX
				mainGameController.cars[i].rotationType = 3
				mainGameController.cars[i].lastStep = true
			}

			if mainGameController.cars[i].speedX == 0 && mainGameController.cars[i].speedY >= 0.0 && mainGameController.cars[i].rotationType == 3 && mainGameController.cars[i].positionY >= 20 && mainGameController.cars[i].routeNumber == 2 {
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = 0
				mainGameController.cars[i].speedY = 0
				mainGameController.cars[i].status = 1
			}
			//Route3

			//Route4
			//Decisions made in base of the position, rotation and destination of car
			if mainGameController.cars[i].speedX == 0 && mainGameController.cars[i].speedY > 0.0 && mainGameController.cars[i].rotationType == 3 && mainGameController.cars[i].positionY >= 15 && mainGameController.cars[i].routeNumber == 3 {
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = -mainGameController.cars[i].startingSpeedY
				mainGameController.cars[i].speedY = 0
				mainGameController.cars[i].rotationType = 2
			}
			if mainGameController.cars[i].speedX < 0 && mainGameController.cars[i].speedY == 0.0 && mainGameController.cars[i].rotationType == 2 && mainGameController.cars[i].positionX <= -235.00 && mainGameController.cars[i].routeNumber == 3 {
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = 0
				mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedY
				mainGameController.cars[i].rotationType = 3
			}
			if mainGameController.cars[i].speedX == 0 && mainGameController.cars[i].speedY >= 0.0 && mainGameController.cars[i].rotationType == 3 && mainGameController.cars[i].positionY >= 365.00 && mainGameController.cars[i].routeNumber == 3 {
				//mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				mainGameController.cars[i].speedX = 0
				mainGameController.cars[i].speedY = 0
				mainGameController.cars[i].status = 1
			}
			//Route4

			//Light1 -40 -40
			//Light2 200 40
			//Light3 -40 40
			//Light4 -200 -40
			//Light5 40 -200
			//Light6 -40 -40
			//Light6 -40 -40

		
			//Stop in traffic lights
			//route1
			if mainGameController.cars[i].routeNumber == 0 && mainGameController.cars[i].positionY <= -50 && mainGameController.cars[i].positionY >= -54 && mainGameController.cars[i].rotationType == 1{
				if mainGameController.semaphores[3].color == 1 {
					mainGameController.cars[i].speedY = 0
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 0 && k>i && k != i{
							mainGameController.cars[k].alreadyStopped = true
						}
					}
				} else {
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 0 {
							mainGameController.cars[k].alreadyStopped = false
						}
					}
					if mainGameController.cars[i].startingSpeedY == 0 {
						mainGameController.cars[i].speedY = -mainGameController.cars[i].startingSpeedX
					} else {
						mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedY
					}
				}
			}

			if mainGameController.cars[i].routeNumber == 0 && mainGameController.cars[i].positionX >= 245 && mainGameController.cars[i].positionX <= 250 {
				if mainGameController.semaphores[6].color == 1 {
					mainGameController.cars[i].speedX = 0
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 0 && k>i && k != i{
							mainGameController.cars[k].alreadyStopped = true
						}
					}
				} else {
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 0 {
							mainGameController.cars[k].alreadyStopped = false
						}
					}
					if mainGameController.cars[i].startingSpeedX == 0 {
						mainGameController.cars[i].speedX = -mainGameController.cars[i].startingSpeedY
					} else {
						mainGameController.cars[i].speedX = mainGameController.cars[i].startingSpeedX
					}
				}
			}
			//

			//route2
			if mainGameController.cars[i].routeNumber == 1 && mainGameController.cars[i].positionY <= -95 && mainGameController.cars[i].positionY >= -105 {
				if mainGameController.semaphores[2].color == 1 {
					mainGameController.cars[i].speedY = 0
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 1 && k>i && k != i{
							mainGameController.cars[k].alreadyStopped = true
						}
					}
				} else {
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 1 {
							mainGameController.cars[k].alreadyStopped = false
						}
					}
					if mainGameController.cars[i].startingSpeedX == 0 {
						mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedY
					} else {
						mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedY
					}
				}
			}

			if mainGameController.cars[i].routeNumber == 1 && mainGameController.cars[i].positionX >= 100 && mainGameController.cars[i].positionX <= 110 {
				if mainGameController.semaphores[4].color == 1 {
					mainGameController.cars[i].speedX = 0
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 1 && k>i && k != i{
							mainGameController.cars[k].alreadyStopped = true
						}
					}
				} else {
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 1 {
							mainGameController.cars[k].alreadyStopped = false
						}
					}
					if mainGameController.cars[i].startingSpeedY == 0 {
						mainGameController.cars[i].speedX = mainGameController.cars[i].startingSpeedX
					} else {
						mainGameController.cars[i].speedX = -mainGameController.cars[i].startingSpeedY
					}
				}
			}
			//

		
			//route3
			if mainGameController.cars[i].routeNumber == 2 && mainGameController.cars[i].positionY <= -50 && mainGameController.cars[i].positionY >= -60 && mainGameController.cars[i].rotationType == 1{
				if mainGameController.semaphores[1].color == 1 {
					mainGameController.cars[i].speedY = 0
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 2 && k>i && k != i{
							mainGameController.cars[k].alreadyStopped = true
						}
					}
				} else {
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 2 {
							mainGameController.cars[k].alreadyStopped = false
						}
					}
					if mainGameController.cars[i].startingSpeedX == 0 {
						mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedY
					} else {
						mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedX
					}
				}
			}
			if mainGameController.cars[i].routeNumber == 2 && mainGameController.cars[i].positionX <= -140 && mainGameController.cars[i].positionX >= -150 {
				if mainGameController.semaphores[5].color == 1 {
					mainGameController.cars[i].speedX = 0
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 2 && k>i && k != i{
							mainGameController.cars[k].alreadyStopped = true
						}
					}
				} else {
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 2 {
							mainGameController.cars[k].alreadyStopped = false
						}
					}
					if mainGameController.cars[i].startingSpeedY == 0 {
						mainGameController.cars[i].speedX = mainGameController.cars[i].startingSpeedX
					} else {
						mainGameController.cars[i].speedX = mainGameController.cars[i].startingSpeedY
					}
				}
			}

			//route4
			if mainGameController.cars[i].routeNumber == 3 && mainGameController.cars[i].positionX <= -65 && mainGameController.cars[i].positionX >= -75 {
				if mainGameController.semaphores[4].color == 1 {
					mainGameController.cars[i].speedX = 0
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 3 && k>i && k != i{
							mainGameController.cars[k].alreadyStopped = true
						}
					}
				} else {
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 3 {
							mainGameController.cars[k].alreadyStopped = false
						}
					}
					if mainGameController.cars[i].startingSpeedY == 0 {
						mainGameController.cars[i].speedX = mainGameController.cars[i].startingSpeedX
					} else {
						mainGameController.cars[i].speedX = -mainGameController.cars[i].startingSpeedY
					}
				}
			}

			if mainGameController.cars[i].routeNumber == 3 && mainGameController.cars[i].positionY >= 135 && mainGameController.cars[i].positionY <= 145 {
				if mainGameController.semaphores[0].color == 1 {
					mainGameController.cars[i].speedY = 0
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 3 && k>i && k != i{
							mainGameController.cars[k].alreadyStopped = true
						}
					}
				} else {
					for k := 0; k < mainGameController.numberOfCars; k++ {
						if mainGameController.cars[k].routeNumber == 3 {
							mainGameController.cars[k].alreadyStopped = false
						}
					}
					if mainGameController.cars[i].startingSpeedX == 0 {
						mainGameController.cars[i].speedY = mainGameController.cars[i].startingSpeedY
					} else {
						mainGameController.cars[i].speedY = -mainGameController.cars[i].startingSpeedX
					}
				}
			}
			//
			//Stop in traffic lights

			screen.DrawImage(mainGameController.cars[i].image, &mainGameController.cars[i].options)
		}
	}

	//This part draws the text of the trafic lights
	for i := 0; i < mainGameController.numberOfSemaphores; i++ {
		if mainGameController.semaphores[i].color == 1 {
			text.Draw(screen, "Rojo", normalFont, mainGameController.traficLightControllers[i].positionX , mainGameController.traficLightControllers[i].positionY-7, color.White)
		}else if mainGameController.semaphores[i].color == 0 {
			text.Draw(screen, "Verde", normalFont, mainGameController.traficLightControllers[i].positionX , mainGameController.traficLightControllers[i].positionY-7, color.White)
		}	
		screen.DrawImage(mainGameController.semaphores[i].image,mainGameController.semaphores[i].options)
	}

	return nil
}

func main() {
	//Initialize seed for real random numbers
	rand.Seed(time.Now().UnixNano())
	//Status of the game
	//0 Game not started
	//1 Game started
	//2 Game finished
	gameStarted = 0
	//Initialization of game controller
	mainGameController = gameController{
		numberOfCars:       7,
		numberOfSemaphores: 7,
		screenWidth: 620,	
		screenHeight: 400,
	}

	//Initialization of starting points and ending points for the routes
	/*
		Routes
		1-Left to right
		2-Bottom to top
		3-right to bottom
		4-Up to bottom
	*/
	positionX1 := 10.0
	positionY1 := 295.0

	positionX2 := 325.0
	positionY2 := 360.0

	positionX3 := 520.0
	positionY3 := 340.0

	positionX4 := 530.0
	positionY4 := 30.0

	//Ending postions
	epositionX1 := 530.0
	epositionY1 := 30.0

	epositionX2 := 530.0
	epositionY2 := 30.0

	epositionX3 := 310.0
	epositionY3 := 360.0

	epositionX4 := 10.0
	epositionY4 := 280.0

	mainGameController.startingPositionsX = append(mainGameController.startingPositionsX,positionX1)
	mainGameController.startingPositionsY = append(mainGameController.startingPositionsY,positionY1)

	mainGameController.startingPositionsX = append(mainGameController.startingPositionsX,positionX2)
	mainGameController.startingPositionsY = append(mainGameController.startingPositionsY,positionY2)

	mainGameController.startingPositionsX = append(mainGameController.startingPositionsX,positionX3)
	mainGameController.startingPositionsY = append(mainGameController.startingPositionsY,positionY3)

	mainGameController.startingPositionsX = append(mainGameController.startingPositionsX,positionX4)
	mainGameController.startingPositionsY = append(mainGameController.startingPositionsY,positionY4)

	//Ending positions
	mainGameController.endingPositionsX = append(mainGameController.endingPositionsX,epositionX1)
	mainGameController.endingPositionsY = append(mainGameController.endingPositionsY,epositionY1)

	mainGameController.endingPositionsX = append(mainGameController.endingPositionsX,epositionX2)
	mainGameController.endingPositionsY = append(mainGameController.endingPositionsY,epositionY2)

	mainGameController.endingPositionsX = append(mainGameController.endingPositionsX,epositionX3)
	mainGameController.endingPositionsY = append(mainGameController.endingPositionsY,epositionY3)

	mainGameController.endingPositionsX = append(mainGameController.endingPositionsX,epositionX4)
	mainGameController.endingPositionsY = append(mainGameController.endingPositionsY,epositionY4)


	//Options opened flag
	optionsOpen := false
	for gameStarted != 3 {
		fmt.Println("Welcome to trafic simulation, select the option you desire")
		fmt.Println("1: Start Game")
		fmt.Println("2: Change Options")
		fmt.Println("3: Exit")
		fmt.Print("Number of the option: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		option := input.Text()
		if option == "1" {
			tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
			if err != nil {
				log.Fatal(err)
			}
			normalFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
				Size:    12,
				DPI:     dpi,
				Hinting: font.HintingFull,
			})
			for i := 0; i < mainGameController.numberOfCars; i++ {
				//Default values for cars
				randomPosition := (rand.Intn(4) + 0)
				currentImage, _, err := ebitenutil.NewImageFromFile("roads/car.png", ebiten.FilterDefault)
				if err != nil {
					log.Fatal(err)
				}
				tempRoute := route{
					startX: mainGameController.startingPositionsX[randomPosition],
					startY: mainGameController.startingPositionsY[randomPosition],
					destinationX: mainGameController.endingPositionsX[randomPosition],
					destinationY: mainGameController.endingPositionsY[randomPosition],
				}
				tempSpeed := (1.5 + rand.Float64() * (1.8-1.5))
				tempCar := car{
					//Random speed from 1 to 10
					speedX: tempSpeed,
					//speedY: (rand.Float64() * (2.5-.8)),
					speedY: 0,
					startingSpeedX: tempSpeed,
					startingSpeedY: 0,
					status: 0,
					rotationType: 0,
					carRoute: tempRoute,
					routeNumber: randomPosition,
					//Image of the car
					image: currentImage,
					positionX: 0,
					positionY: 0,
					//Draw Options of the car
					options: ebiten.DrawImageOptions{},
					lastStep: false,
					alreadyStopped: false,
					startTime: 2 + (i+2),
				}
				mainGameController.cars = append(mainGameController.cars, tempCar)
				if randomPosition == 0{
					//Rotation to right
					mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*360))
					mainGameController.cars[i].rotationType = 0;
				}else if randomPosition == 1{
					//Rotation top
					mainGameController.cars[i].speedX = 0
					mainGameController.cars[i].speedY = -tempSpeed
					mainGameController.cars[i].startingSpeedX = 0
					mainGameController.cars[i].startingSpeedY = -tempSpeed
					mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*180))
					mainGameController.cars[i].rotationType = 1;
				}else if randomPosition == 2{
					//Rotation left
					mainGameController.cars[i].speedX = -tempSpeed
					mainGameController.cars[i].speedY = 0
					mainGameController.cars[i].startingSpeedX = -tempSpeed
					mainGameController.cars[i].startingSpeedY = 0
					mainGameController.cars[i].rotationType = 2;
				}else if randomPosition == 3{
					//Rotation Bottom
					mainGameController.cars[i].speedX = 0
					mainGameController.cars[i].speedY = tempSpeed
					mainGameController.cars[i].startingSpeedX = 0
					mainGameController.cars[i].startingSpeedY = tempSpeed
					mainGameController.cars[i].rotationType = 3;
					mainGameController.cars[i].options.GeoM.Rotate(float64((math.Pi / 360)*540))
				}
				mainGameController.cars[i].options.GeoM.Translate(mainGameController.cars[i].carRoute.startX, mainGameController.cars[i].carRoute.startY)
				go carBehavior(i)
			}
			for i := 0; i < 7; i++ {
				//Default values for semaphores
				tempSemaphore := semaphore{
					//Random color from 1 to 0
					color:   (rand.Intn(2) + 0),
					counter: 0,
				}
				mainGameController.semaphores = append(mainGameController.semaphores, tempSemaphore)
				go semaphoreBehavior(i)
			}
			gameStarted = 3
			if err := ebiten.Run(update, mainGameController.screenWidth, mainGameController.screenHeight, 2, "Traffic Simulator"); err != nil {
				log.Fatal(err)
			}
			fmt.Println(" ")
		} else if option == "2" {
			optionsOpen = true
			fmt.Println(" ")
			for optionsOpen == true {
				fmt.Println("Options menu, select the option you desire to change")
				fmt.Println("1: Select number of cars")
				fmt.Println("2: Select number of semaphores")
				fmt.Println("3: Return to main menu")
				fmt.Print("Number of the option: ")
				input := bufio.NewScanner(os.Stdin)
				input.Scan()
				specificOption := input.Text()
				fmt.Println(" ")
				if specificOption == "1" {
					fmt.Print("Select the number of cars: ")
					input := bufio.NewScanner(os.Stdin)
					input.Scan()
					specificOptionCars := input.Text()
					cars, _ := strconv.Atoi(specificOptionCars)
					if(cars > 8){
						fmt.Println("Please select a valid option")
						fmt.Println(" ")
						fmt.Println(" ")
					}else{
						mainGameController.numberOfCars = cars
						fmt.Println("Number of cars changed")
						fmt.Println(" ")
					}
				} else if specificOption == "2" {
					fmt.Print("Select the number of traffic lights: ")
					input := bufio.NewScanner(os.Stdin)
					input.Scan()
					specificOptionSemaphores,_ := strconv.Atoi(input.Text())
					if(specificOptionSemaphores > 7){
						fmt.Println("Please select a valid option")
						fmt.Println(" ")
						fmt.Println(" ")
					}else{
						semaphores := specificOptionSemaphores
						mainGameController.numberOfSemaphores = semaphores
						fmt.Println("Number of semaphores changed")
						fmt.Println(" ")
					}
				} else if specificOption == "3" {
					optionsOpen = false
					fmt.Println(" ")
				} else {
					fmt.Println("Please select a valid option")
					fmt.Println(" ")
					fmt.Println(" ")
				}
			}
		} else if option == "3" {
			fmt.Println("Thanks for playing Traffic Simulator")
			os.Exit(1)
		} else if option == "4" {

		} else {
			fmt.Print("Please select a valid option")
			fmt.Println(" ")
		}
	}

}
