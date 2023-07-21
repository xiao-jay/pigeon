package festival

//import (
//	"fmt"
//	"github.com/PuloV/ics-golang"
//)
//
//func init() {
//	fmt.Println("1")
//	parser := ics.New()
//	parserChan := parser.GetInputChan()
//	parserChan <- "https://rili.callmekeji.com/calendar/keji_calendar.ics"
//
//	outputChan := parser.GetOutputChan()
//	//  print events
//	go func() {
//		for event := range outputChan {
//			fmt.Println(event.GetImportedID(), event.GetSummary())
//		}
//	}()
//
//	// wait to kill the main goroute
//	parser.Wait()
//}
