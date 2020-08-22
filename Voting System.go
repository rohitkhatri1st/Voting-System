package main

import (
	"encoding/json"
	"fmt"
	"greet"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"

	"gopkg.in/gookit/color.v1"
)

type Person struct {
	First    string
	Last     string
	Id       int32
	Password []uint8
	Voted    map[string]bool
	//Admin_code string
}
type Candidate struct {
	First  string
	Last   string
	Symbol string
	Count  int
	//Id     int
	//Admin_code string
}

func main() {
	var IsNew int = 0
	var FirstName, LastName string
	var ID int32
	var Pass []uint8
	var IsLoggedIn bool
	var p1 Person
	dologin := false
	m := map[int32]Person{}
	c := map[string]Candidate{}
	//var people []Person
	greet.HomeScreen()
	defer func() {
		time.Sleep(3 * time.Second)
	}()
	_, err := fmt.Scan(&IsNew)
	for err != nil || (IsNew != 2 && IsNew != 1 && IsNew != 3) {
		if err != nil {
			color.Red.Println(err)
		}
		color.Red.Println(" Please input 1 for new user, 2 for old user and 3 if you are an admin.")
		_, err = fmt.Scan(&IsNew)
	}
	Cls()
	f, err := os.OpenFile("myfile.data", os.O_RDONLY|os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadFile("myfile.data")
	os.MkdirAll("data", 0777)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(data, &m)
	defer f.Close()
	//fmt.Println("Candidates list as c is", c)
	if IsNew == 1 {
		color.Cyan.Println("\t\tHello User... Please Enter your details to Proceed...")
		color.Yellow.Print("Enter your Id Number:\t\t")
		ok := true
		for ok == true {
			_, err = fmt.Scan(&ID)
			if err != nil {
				color.Red.Println(err)
				continue
			} else if _, ok := m[ID]; ok {
				color.Red.Println("Entered Id number is already in use. Please try with a new Id number.")
				continue
			}
			ok = false
		}
		color.Yellow.Print("First Name:\t\t")
		fmt.Scan(&FirstName)
		color.Yellow.Print("Last Name:\t\t")
		fmt.Scan(&LastName)
		color.Yellow.Print("Enter a Password:\t")
		fmt.Scan(&Pass)
		HashedPass, _ := bcrypt.GenerateFromPassword([]byte(Pass), bcrypt.MinCost)
		voted := map[string]bool{}
		p1 = Person{FirstName, LastName, ID, HashedPass, voted}
		m[ID] = p1
		bs, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
		}
		f.Write(bs)
		color.Cyan.Print("Hello ", p1.First, " ", p1.Last, "...\n")
		color.LightMagenta.Print("Press 1 to login or any other key to exit program:\t")
		fmt.Scan(&dologin)
		if dologin == true {
			Cls()
		} else {
			fmt.Println("Exiting program...")
			time.Sleep(3 * time.Second)
			os.Exit(1)
		}
		//greet.NewUserScreen()
	}
	if IsNew == 2 || dologin == true {
		color.Cyan.Println("\t\tHello User... Please Login to Proceed.")
		color.Yellow.Print("Enter your Id Number:\t\t")
		_, err = fmt.Scan(&ID)
		errcheck(err)
		for err != nil {
			_, err = fmt.Scan(&ID)
			errcheck(err)
		}
		color.Yellow.Print("Enter a Password:\t")
		//fmt.Scan(&Pass)
		Pass, err = terminal.ReadPassword(int(os.Stdin.Fd()))
		//errcheck(err)
		json.Unmarshal(data, &m)
		err = bcrypt.CompareHashAndPassword(m[ID].Password, []byte(Pass))
		if err != nil {
			color.Red.Println("Entered combination of ID and Password does not match any combination in our database.")
		} else {
			IsLoggedIn = true
		}
		if IsLoggedIn == true {
			//UserLoggedinpage
			fmt.Println("Login Successful...")
			time.Sleep(time.Second)
			Cls()
			color.Cyan.Print("\n\t\tWelcome ", m[ID].First, " ", m[ID].Last, "...\n")
			ElectionListBS, err := ioutil.ReadFile("data/Election.List")
			if err != nil {
				fmt.Println("No ongoing elections currently.")
				Exitkey()
			}
			ListOfElection := make([]map[string]bool, 10)
			json.Unmarshal(ElectionListBS, &ListOfElection)
			color.Cyan.Println("List of ongoing election is as follows:")
			var ElectionNum int
			for i, v := range ListOfElection {
				color.Yellow.Print(i+1, ":\t")
				for k, val := range v {
					color.Yellow.Print(k)
					color.Blue.Println("\tResult Published :", val)
				}
			}
			color.LightMagenta.Println("Choose from above to vote or see result of an election.")
			_, err = fmt.Scan(&ElectionNum)
			for err != nil && ElectionNum < len(ListOfElection) {
				color.Red.Println("Invalid Entry... Please Try again.")
				fmt.Scan(&ElectionNum)
			}
			Cls()
			var ElectionName string
			for k := range ListOfElection[ElectionNum-1] {
				ElectionName = k
			}
			cndidtdata, err := ioutil.ReadFile(ElectionName)
			errcheck(err)
			json.Unmarshal(cndidtdata, &c)
			if ListOfElection[ElectionNum-1][ElectionName] == true {
				color.LightMagenta.Println("Election is over. Below are the results.")
				maxCount := 0
				maxCountCandidateFirstName := ""
				maxCountCandidateLastName := ""
				maxCountCandidateSymbol := ""
				for key, value := range c {
					if maxCount < value.Count {
						maxCount = value.Count
						maxCountCandidateFirstName = value.First
						maxCountCandidateLastName = value.Last
						maxCountCandidateSymbol = key
					}
				}
				fmt.Println(maxCountCandidateFirstName, maxCountCandidateLastName, "with symbol", maxCountCandidateSymbol, "has won the election with total number of votes=\t", maxCount)
				var exitkey int
				fmt.Println("Enter any key to exit")
				fmt.Scan(&exitkey)
				os.Exit(1)
			} else {
				//fmt.Println(m[ID].Voted[ElectionName])
				if m[ID].Voted[ElectionName] == true {
					color.BgRed.Println("You have already voted for", ElectionName, "Election.")
					color.BgRed.Println("For security reasons there is no data regarding who you voted to.")
					Exitkey()
				} else {
					color.LightMagenta.Println("Candidates for election of", ElectionName, "are:")
					color.BgGreen.Print("Symbol\t\t\t")
					color.BgBlue.Println("Name\t\t\t")
					// color.BgRed.Print("Votes\n")
					for key, value := range c {
						color.Green.Print(key)
						color.Blue.Println("\t\t\t", value.First, value.Last)
						// color.BgRed.Println("\t\t\t", value.Count)
					}
					color.Red.Println("NOTE:\tOnce voted can't be undone. Please proceed if sure.\n Press 1 to Proceed, Otherwise input any other key to logout and exit.")
					surity := 0
					fmt.Scan(&surity)
					if surity != 1 {
						fmt.Println("Ok Bye... Exiting...")
						os.Exit(1)
					}
					color.LightMagenta.Println("Enter a symbol to vote for.")
					votedfor := ""
					fmt.Scan(&votedfor)
					for _, ok := c[votedfor]; !ok; {
						color.Red.Println("Not a valid symbol. Please Try again.")
						fmt.Scan(&votedfor)
						_, ok = c[votedfor]
					}
					// c1 := Candidate{CandidateFirstName, CandidateLastName, CandidateSymbol, 0}
					// c[CandidateSymbol] = c1
					c[votedfor] = Candidate{c[votedfor].First, c[votedfor].Last, c[votedfor].Symbol, c[votedfor].Count + 1}
					m[ID].Voted[ElectionName] = true
					mBs, err := json.Marshal(m)
					errcheck(err)
					ioutil.WriteFile("myfile.data", mBs, 0777)
					cndidtdataupdated, err := json.Marshal(c)
					errcheck(err)
					err = ioutil.WriteFile(ElectionName, cndidtdataupdated, 0777)
					errcheck(err)
					fmt.Println("Voting Successful... Exiting...")
					Exitkey()
				}
			}
			// color.Cyan.Print("Result of ", ElectionName, "(Published) is as follows.\n")
			// // } else {
			// // 	color.Cyan.Print("Result of ", ElectionName, "(UnPublished) is as follows.\n")

			// // }
			// color.BgGreen.Print("Symbol\t\t\t")
			// color.BgBlue.Print("Name\t\t\t")
			// //color.BgRed.Print("Votes\n")
			// for key, value := range c {
			// 	color.Green.Print(key)
			// 	color.Blue.Println("\t\t\t", value.First, " ", value.Last)
			// 	//color.BgRed.Println("\t\t\t", value.Count)
			// }

			//color.Yellow.Println("1.\tVote\n2.\tUpdate Account Details\n3.\tWho won\n\t\t Or any other key to logout and exit.")
			// functionality := 0
			// fmt.Scan(&functionality)
			// switch functionality {
			// case 1:
			// 	Cls()
			// 	votedfor := 0
			// 	color.Red.Println("NOTE:\tOnce voted can't be undone. Please proceed if sure. Otherwise input any other key to logout and exit.")
			// 	fmt.Scan(&votedfor)
			// 	if votedfor == 1 || votedfor == 2 || votedfor == 3 || votedfor == 4 {

			// 		fmt.Println("You Voted for", votedfor)
			// 	} else {
			// 		IsLoggedIn = false
			// 		fmt.Println("Ok Bye...")
			// 		os.Exit(1)
			// 	}
			// case 2:
			// 	fmt.Println("You chose update account details")
			// case 3:
			// 	fmt.Println("You chose who won.")
			// default:
			// 	IsLoggedIn = false
			// 	fmt.Println("Ok bye...")
			// }
		} else {
			fmt.Println("Something seems wrong... :(\nI am exiting!!")
		}
	} else {
		//Cls()
		color.LightMagenta.Println("\tEnter the admin code:")
		Admincode := ""
		fmt.Scan(&Admincode)
		if Admincode == "Iamatrustedadmin" {
			color.Cyan.Println("\t\t\tWelcome Admin...")
			color.LightMagenta.Println("Choose:")
			color.Yellow.Println("1.\tConduct an election\n2.\tShow or Update Ongoing Election\n")
			functionality := 0
			_, err := fmt.Scan(&functionality)
			errcheck(err)
			switch functionality {
			case 1:
				color.LightMagenta.Println("Why is the Election being held?")
				ElectionPost := map[string]bool{}
				ElectionListBS, err := ioutil.ReadFile("data/Election.List")
				//errcheck(err)
				ListOfElection := []map[string]bool{}
				json.Unmarshal(ElectionListBS, &ListOfElection)
				//ioutil.ReadFile("ListOfElection.data")
				var ElectionPostvar string
				_, err = fmt.Scan(&ElectionPostvar)
				errcheck(err)
				ElectionPost[ElectionPostvar] = false
				ListOfElection = append(ListOfElection, ElectionPost)
				ListOfElectionBS, err := json.Marshal(ListOfElection)
				errcheck(err)
				err = ioutil.WriteFile("data/Election.List", ListOfElectionBS, 0777)
				color.LightMagenta.Println("Number of candidates?")
				var NumofCandidates int
				_, err = fmt.Scan(&NumofCandidates)
				for err != nil {
					color.Red.Println(err)
					_, err = fmt.Scan(&NumofCandidates)
				}
				for i := 1; i <= NumofCandidates; i++ {
					color.LightMagenta.Println("Input details of candiate", i, ":")
					CandidateInput(c, ElectionPostvar)
				}
				color.Green.Println("Election List updated successfully. Enter any key to exit program.")
				var exitstatus bool
				fmt.Scan(&exitstatus)
				os.Exit(1)
				// fmt.Println("List of candidates is as follows:\n")
				// for k, v := range c {
				// 	color.Green.Println(k, v)
				// }

			case 2:
				ElectionListBS, err := ioutil.ReadFile("data/Election.List")
				if err != nil {
					fmt.Println("No ongoing elections currently.")
				}
				ListOfElection := make([]map[string]bool, 10)
				json.Unmarshal(ElectionListBS, &ListOfElection)
				color.Cyan.Println("List of ongoing election is as follows:")
				var ElectionNum int
				for i, v := range ListOfElection {
					color.Yellow.Print(i+1, ":\t")
					for k, val := range v {
						color.Yellow.Print(k)
						color.Blue.Println("\tPublished status:", val)
					}
				}
				color.LightMagenta.Println("Choose from above to check, publish or unpublish result of any election.")
				fmt.Scan(&ElectionNum)
				Cls()
				var ElectionName string
				for k := range ListOfElection[ElectionNum-1] {
					ElectionName = k
				}
				if ListOfElection[ElectionNum-1][ElectionName] == true {
					color.Cyan.Print("Result of ", ElectionName, "(Published) is as follows.\n")
				} else {
					color.Cyan.Print("Result of ", ElectionName, "(UnPublished) is as follows.\n")

				}
				cndidtdata, err := ioutil.ReadFile(ElectionName)
				errcheck(err)
				json.Unmarshal(cndidtdata, &c)
				color.BgGreen.Print("Symbol\t\t\t")
				color.BgBlue.Print("Name\t\t\t")
				color.BgRed.Print("Votes\n")
				for key, value := range c {
					color.Green.Print(key)
					color.Blue.Print("\t\t\t", value.First, " ", value.Last)
					color.BgRed.Println("\t\t\t", value.Count)
				}
				isPublished := 0
				color.Yellow.Println("Enter 1 to publish/unpublish result or any other key to exit.")
				fmt.Scan(&isPublished)
				if isPublished == 1 {
					if ListOfElection[ElectionNum-1][ElectionName] == true {
						ListOfElection[ElectionNum-1][ElectionName] = false
					} else {
						ListOfElection[ElectionNum-1][ElectionName] = true
					}
					ListOfElectionBS, err := json.Marshal(ListOfElection)
					err = ioutil.WriteFile("data/Election.List", ListOfElectionBS, 0777)
					if err != nil {
						color.Println(err)
					} else if ListOfElection[ElectionNum-1][ElectionName] == true {
						fmt.Println("Result Published Successfully. Please Restart.")
					} else {
						fmt.Println("Result unpublished Successfully. Please Restart.")
					}
				} else {
					os.Exit(1)
				}
			}
		} else {
			fmt.Println("Bahh!! :(\tSomething seems wrong. I am exiting")
			os.Exit(0)
		}
	}
}
func CandidateInput(c map[string]Candidate, ElectionPostvar string) {
	cndidt, err := os.OpenFile(ElectionPostvar, os.O_CREATE|os.O_RDONLY|os.O_WRONLY, 0700)
	errcheck(err)
	cndidtdata, err := ioutil.ReadFile(ElectionPostvar)
	errcheck(err)
	json.Unmarshal(cndidtdata, &c)
	defer cndidt.Close()
	// color.Yellow.Print("ID number :\t\t")
	// var ID int
	// ok := true
	// for ok == true {
	// 	_, err = fmt.Scan(&ID)
	// 	if err != nil {
	// 		color.Red.Println(err)
	// 		continue
	// 	} else if _, ok := c[ID]; ok {
	// 		color.Red.Println("Entered ID is already in use. Please try with a new ID.")
	// 		continue
	// 	}
	// 	ok = false
	// }
	color.Yellow.Print("First Name :\t\t")
	var CandidateFirstName string
	fmt.Scan(&CandidateFirstName)
	color.Yellow.Print("Last Name :\t\t")
	var CandidateLastName string
	fmt.Scan(&CandidateLastName)
	color.Yellow.Print("Symbol :\t\t")
	var CandidateSymbol string
	ok := true
	for ok == true {
		_, err = fmt.Scan(&CandidateSymbol)
		if err != nil {
			color.Red.Println(err)
			continue
		} else if _, ok := c[CandidateSymbol]; ok {
			color.Red.Println("Entered symbol is already in use. Please try with a new symbol.")
			continue
		}
		ok = false
	}
	c1 := Candidate{CandidateFirstName, CandidateLastName, CandidateSymbol, 0}
	c[CandidateSymbol] = c1
	Candidatebs, err := json.Marshal(c)
	errcheck(err)
	_, err = cndidt.Write(Candidatebs)
	errcheck(err)
}
func Cls() {
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
}
func errcheck(err error) {
	if err != nil {
		color.Red.Println(err)
	}
}
func Exitkey() {
	fmt.Println("Enter any key to exit program.")
	var key int
	fmt.Scan(&key)
	os.Exit(1)
}
