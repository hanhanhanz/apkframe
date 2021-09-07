package main
//import "log"
import "fmt"
import "os"
import "io/ioutil"
import "path/filepath"
import "bufio"
import "strings"
import "strconv"
import "flag"
import "os/exec"
//import "sync"
import "encoding/json"
//import "time"


type Fingerprints struct {
	Fingerprints []Fingerprint `json:"fingerprints"`
}

type Fingerprint struct{
	Name 		string 			`json:"name"`
	File 		[]string 		`json:"file"`
	Dir 		[]string 		`json:"dir"`
	Findstring	[]Findstring 	`json:"findstring"`
}

type Findstring struct{
	Dirtofind 	string 			`json:"dirtofind"`
	File 		[]string 		`json:"file"`
	Str  		string 			`json:"Str"` 
}


func libcheck(path string) string {
	if strings.Contains(path, "§") {
    	paths := strings.Split(path, "§")
		path = paths[0]
		//fmt.Println(path)
		//nge-ls di golang
		//files, err := ioutil.ReadDir(path)
		files, err := ioutil.ReadDir(path)
		
		if err != nil {
	        if !strings.Contains(err.Error(), "no such file or directory") {
			 	fmt.Println(err)
	    	}
	    	return path
		}
	 	//print directory
	    /*for _, f := range files {
	            fmt.Println(f.Name())
	    }*/
	    if len(files) > 0 {
	    	path = paths[0] + files[0].Name() +  paths[1]	
	    	return path
	    }
	    

    }
    return path
}
    

func dirf(path string) (bool, error) {
    infofile, err := os.Stat(path)
    if err == nil { 
    	if infofile.IsDir() {
    		return true, nil 
    	}	
    }
    if err != nil { return false, err }
    return false, err
}

func filef(path string) (bool, error) {
    


    path = libcheck(path)

    infofile, err := os.Stat(path)
    
    if err == nil { 
    	if !(infofile.IsDir()) {
    		return true, nil 
    	}	
    }
    if err != nil { return false, err }
    return false, err
}

//func filef

func openandfind(path string, seed string) (string,int,error) {
	f, err := os.Open(path)
	if err != nil {
	    return "",0, err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
	    if strings.Contains(scanner.Text(), seed) {
	        //fmt.Println(path)
			//fmt.Println(strconv.Itoa(line))
	        return path, line, nil
	    }

	    line = line + 1
	}

	if err := scanner.Err(); err != nil {
	    
	    return "",0, err
	}
	
	return "",0, err
}

func main() {

	apk := ""
	var apkt bool
	flag.StringVar(&(apk),"d","","specify output directory of apktool")
	flag.BoolVar(&(apkt),"a",false,"decompile apk with ")
	flag.Parse()
	if apk == "" {
        panic("no apk dude")
    }

	if apkt == true {
		arg0 := "apktool"
		arg1 := "d"
		arg2 := apk
		arg3 := "-o"
		arg4 := fmt.Sprintf("ss-"+filepath.Base(apk))
		if _, err := os.Stat("ss-"+filepath.Base(apk)); os.IsNotExist(err) {
			cmd := exec.Command(arg0,arg1,arg2,arg3,arg4)
	    	err2 := cmd.Run()
	    	if err2 == nil {
	    		apk =  "ss-"+ filepath.Base(apk)
	        	
	    	}
		} else {
			fmt.Println("[!] using cached decompiler")
		}	
	}


	/*currdir, _ := os.Getwd()
	_, err := os.Stat(currdir + "/" + apk)
	if err != nil { 
    	_, err := os.Stat(apk)
    	if err != nil { 
    		panic("missing dir/apk")
    	}
    	
    } else {
    	apk = currdir + "/" + apk
    }*/



	jsonFile, err := os.Open("signature.json")
	if err != nil{
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var f Fingerprints
	json.Unmarshal(byteValue, &f)
	//fmt.Println(len(f.Fingerprints))
	hit := map[string]int{}
	for i := 0; i < len(f.Fingerprints); i++ {
		//fmt.Println("fingerprintnya : " + f.Fingerprints[i].Name)
		for j := 0; j < len(f.Fingerprints[i].File); j++ {
			b,err := filef(apk + "/" + f.Fingerprints[i].File[j])
			//if err != nil { fmt.Println(err) }
			if (b == true) && (err == nil) {
				//fmt.Println("signatureMatch = " + f.Fingerprints[i].Name)
				
				//fmt.Println(hit[f.Fingerprints[i].Name])
				
				if hit[f.Fingerprints[i].Name] == 0 {

					hit[f.Fingerprints[i].Name] = 1
				} else {
					hit[f.Fingerprints[i].Name] = hit[f.Fingerprints[i].Name] + 1
				}

			}
		}
		for j := 0; j < len(f.Fingerprints[i].Dir); j++ {
			b,err := dirf(apk + "/" + f.Fingerprints[i].Dir[j])
			//if err != nil { fmt.Println(err) }
			if (b == true) && (err == nil)  {
				//fmt.Println("signatureMatch = " + f.Fingerprints[i].Name)
				if hit[f.Fingerprints[i].Name] == 0 {

					hit[f.Fingerprints[i].Name] = 1
				} else {
					hit[f.Fingerprints[i].Name] = hit[f.Fingerprints[i].Name] + 1
				}
			}
		}

		var dirlist = []string{}
		for j := 0; j < len(f.Fingerprints[i].Findstring); j++ {
			for k := 0; k < len(f.Fingerprints[i].Findstring[j].File); k++ {
				err = filepath.Walk(apk,func(path string, info os.FileInfo, err error) error {
			    if err != nil {
			        return err
			    }
				    if strings.HasSuffix(path,f.Fingerprints[i].Findstring[j].File[k]) {
						//fmt.Println(path)
						dirlist = append(dirlist, path)
				    	return nil	
				    }	
				    return nil	
				})
				if err != nil {
				    fmt.Println(err)
				}
			//fmt.Println(dirlist)
			}
			for k := 0; k < len(dirlist); k++ {		
				//b := findstringf(apk,f.Fingerprints[i].Findstring[j].Dirtofind,f.Fingerprints[i].Findstring[j].Str) //rootdir,assets,framework7
				_,line,err := openandfind(dirlist[k],f.Fingerprints[i].Findstring[j].Str)  //path,line,err := openandfind(dirlist[k],f.Fingerprints[i].Findstring[j].Str)
				if err != nil {
					//fmt.Println("error walking the path")
					fmt.Println(err)
					os.Exit(3)
				}
				if line != 0 {
					//fmt.Println(path +":"+ strconv.Itoa(line))
					if hit[f.Fingerprints[i].Name] == 0 {

						hit[f.Fingerprints[i].Name] = 1
					} else {
						hit[f.Fingerprints[i].Name] = hit[f.Fingerprints[i].Name] + 1
					}
					
				}
				
			}
		
		}
		
	}
	//fmt.Println(hit)
	for key,val := range hit{
		//fmt.Println(apk)
		fmt.Println(strconv.Itoa(val) +" "+ key + " signature found")
		//fmt.Println(key, " : ",val)
	}
}


