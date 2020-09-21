package main
import(
	"os"
	"encoding/csv"
	"log"
	"io"
	//"fmt"
)

type Student struct{
	studentName string
	preferEmployer []string
	match string
}
type Employer struct {
	student []Student
	employerName string
	perferStudent []string
	match string
	offerIndex int
}
//Read files employers and students then return the data inside the file.
func read(eFile string, sFile string) ([]Employer, []Student) {
	//read employer file
	csvFile, _ := os.Open(eFile)
	reader := csv.NewReader(csvFile)
	//read student file
	csvFile2, _ := os.Open(sFile)
	reader2 := csv.NewReader(csvFile2)

	count := 0
	count2 := 0
	var employers []Employer
	var students []Student
	//assign employers to the list
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil{
			log.Fatal(error)
		}
		//create the employer
		employers = append(employers, Employer {
			employerName: line[0],
			match: "None",
			offerIndex: 0,
		})
		//Assign perfer students to the employers
		for i:=0; i< len(line)-1;i++{
			employers[count].perferStudent = append(employers[count].perferStudent, line[i+1])
		}
		count++
	}

	//assign students to the list
	for {
		line2, error := reader2.Read()
		if error == io.EOF {
			break
		} else if error != nil{
			log.Fatal(error)
		}
		//create the employer
		students = append(students, Student{
			studentName: line2[0],
			match: "None",
		})
		//Assign perfer students to the employers
		for m:=0; m < len(line2)-1; m++ {
			students[count2].preferEmployer = append(students[count2].preferEmployer, line2[m+1])
		}
		count2++
	}

	return employers, students
}

//display result function
func displayResult(employers []Employer) []string {
	var displayList []string
	for ii:=0;ii<2;ii++{
		if ii == 0{
			for i := 0 ; i < len(employers); i++ {
				//this will add all of the students name
				displayList = append(displayList, employers[i].match)
			}
		} else{
			for i := 0 ; i < len(employers); i++ {
				//this will add all of the students name
				displayList = append(displayList, employers[ii].student[i].match)
			}
		}
	}
	return displayList
}

func main() {
	//read files and return the values
	employers, students := read("coop_e_3x3.csv", "coop_s_3x3.csv")
	employers2, students2 := read("coop_e_10x10.csv", "coop_s_10x10.csv")
	displayList  := run(employers, students)
	displayList2 := run(employers2, students2)
	createCSV(displayList, "matches_go_3x3.csv")
	createCSV(displayList2, "matches_go_10x10.csv")
}

func run(employers []Employer, students []Student) []string{
	for i := 0; i < len(employers);i++{
		employers[i].student = students
	}	
	stop := false
    for stop == false {
		count := len(employers) * len(employers)
		for i := 0 ; i < len(employers); i++{
			employers = offer(employers[i], employers, i)
		}
		for i:= 0; i < len(employers); i++{
			for m := 0; m < len(employers); m++{
				if employers[i].student[m].match != "None" {
					count--
				}
			}
		}
		if (count == 0){
			stop = true
		}
	}
	displayList := displayResult(employers)
	return displayList
}

func offer(e Employer, employers []Employer, employerI int) []Employer {
	//if e.match = None then it is unmatched
	if e.match == "None"{
		//find the first perfer student name
		//offerIndex keep tract of the students that has not been offer the job.
		studentName := e.perferStudent[e.offerIndex]
		// increase offer Index which means that student index-1 has/had been offered
		e.offerIndex++

		//find the data of the student so it can be passed on to the next function
		for i:=0; i < len(e.student); i++{
			if studentName == e.student[i].studentName{
				//save the student variable that found
				s := e.student[i]
				employers = evaluate(s,e, employers, employerI, i)
				break
			}
		}

		//now we found the most preferred student s on the list the employer
		// has not yet to offered a job
	}
	return employers
}
func evaluate(s Student, e Employer, employers []Employer, employerI int, studentI int) []Employer {
	//if the student has not yet to be match with an employer
	if s.match == "None"{
		//match student with the employer
		for m := 0 ; m < len(employers);m++{
			employers[m].student[studentI].match = e.employerName
		}
		employers[employerI].match = s.studentName
		//if no comarison and sign immediate then return none
	} else {
		//get the index of current and new in employerperfer list of the student.
		current := s.match
		new := e.employerName

		for i := 0 ; i < len(employers); i ++{
			//if current comes first then that means the current is more favarable
			if current == employers[employerI].student[studentI].preferEmployer[i]{
				//reject offer from e and call the offer function again to find the next student
				employers = offer(e, employers, employerI)
				break
			} else if new == employers[employerI].student[studentI].preferEmployer[i]{
				//if new comes first then that means new is more favarable --> perform swap.
				for m:= 0; m < len(employers);m++{
					employers[m].student[studentI].match = e.employerName
				}
				employers[employerI].match = s.studentName
				//to set the current employer match to none, we have to search for him
				for nn := 0; nn < len(employers); nn++{
					if current == employers[nn].employerName{
						employers[nn].match ="None"
					}
				}
				break
			}
		}
	}
	return employers
}

func createCSV (displayList []string, fileName string) {
	var newList []string
	for i := 0; i < len(displayList)/2 ; i++{
		newList = append(newList, displayList[i])
		newList = append(newList, displayList[i + len(displayList)/2])
	}

	display := make ([][]string, len(displayList)/2)
	p:=0;
	for i:=0; i < len(display); i++{
		for m := 0; m < 2; m++{
			display[i] = append(display[i], newList[p])
			p++;
		}
	}

	file,err := os.Create(fileName)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()
	defer file.Close()
    
    writer.WriteAll(display)
}