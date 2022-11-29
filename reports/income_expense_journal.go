package reports

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tadovas/iv-tracker/income"
	"github.com/xuri/excelize/v2"
)

type JournalGenerator struct {
	IncomeRepository income.Repository
}

type JournalGeneratorParams struct {
	name       string
	familyname string
	ssn        string
	address    string
}

func (generator JournalGenerator) GenerateIncomeExpenseJournalAsExcel() ([]byte, error) {
	f := excelize.NewFile()

	yearSummaries, err := generator.IncomeRepository.ListAllYears()
	if err != nil {
		return nil, errors.Errorf("income years list: %v", err)
	}
	for _, yearSummary := range yearSummaries {
		sheetIndex := f.NewSheet(fmt.Sprint(yearSummary.Year))
		f.SetActiveSheet(sheetIndex)
		if err := generator.generateYearSheet(yearSummary, f, fmt.Sprint(yearSummary.Year), 1); err != nil {
			return nil, errors.Errorf("income year gen: %v", err)
		}
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (generator JournalGenerator) generateYearSheet(yearSummary income.YearSummary, excelFile *excelize.File, sheet string, cellOffset int) error {
	incomeList, err := generator.IncomeRepository.ListIncomesByYear(yearSummary.Year)
	if err != nil {
		return errors.Errorf("income list for year: %v", err)
	}
	incomeOffset, err := generateIncomeHeader(excelFile, sheet, cellOffset)
	if err != nil {
		return errors.Errorf("header format: %v", err)
	}
	for i, incomeRecord := range incomeList {
		if err := generateIncomeRow(excelFile, sheet, incomeRecord, i+1, incomeOffset); err != nil {
			return errors.Errorf("income row: %v", err)
		}
		incomeOffset++
	}
	// pajamu suma pagal eilutes
	if err := excelFile.SetCellValue(sheet, asCell("E", incomeOffset), yearSummary.TotalAmount); err != nil {
		return errors.Errorf("sum row: %v", err)
	}
	return nil
}

func generateIncomeRow(excelFile *excelize.File, sheet string, incomeRecord income.Income, incomeNumber int, cellOffset int) error {
	if err := excelFile.SetCellValue(sheet, asCell("A", cellOffset), incomeNumber); err != nil {
		return err
	}
	if err := excelFile.SetCellStr(sheet, asCell("B", cellOffset), incomeRecord.Date.String()); err != nil {
		return err
	}
	if err := excelFile.SetCellValue(sheet, asCell("C", cellOffset), fmt.Sprintf("Pagal kontraktą %v %v", incomeRecord.Origin, incomeRecord.Comment)); err != nil {
		return err
	}
	if err := excelFile.SetCellValue(sheet, asCell("D", cellOffset), "Konsultanto paslaugos pagal kontraktą"); err != nil {
		return err
	}
	if err := excelFile.SetCellValue(sheet, asCell("E", cellOffset), incomeRecord.Amount); err != nil {
		return err
	}
	return nil
}

func generateIncomeHeader(excelFile *excelize.File, sheet string, cellOffset int) (int, error) {

	if err := excelFile.SetCellValue(sheet, asCell("A", cellOffset), "Eilės numeris"); err != nil {
		return 0, err
	}
	if err := excelFile.SetCellValue(sheet, asCell("B", cellOffset), "Data"); err != nil {
		return 0, err
	}
	if err := excelFile.SetCellValue(sheet, asCell("C", cellOffset), "Dokumento data, pavadinimas ir numeris"); err != nil {
		return 0, err
	}
	if err := excelFile.SetCellValue(sheet, asCell("D", cellOffset), "Operacijos turinys"); err != nil {
		return 0, err
	}
	if err := excelFile.SetCellValue(sheet, asCell("E", cellOffset), "Pajamų suma"); err != nil {
		return 0, err
	}

	return cellOffset + 1, nil
}

func asCell(letter string, offset int) string {
	return fmt.Sprintf("%v%v", letter, offset)
}
