package onlineconf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/alldroll/cdb"
	"github.com/stretchr/testify/suite"
)

type testCDBRecord struct {
	key []byte
	val []byte
}
type OCTestSuite struct {
	suite.Suite
	cdbFile   *os.File
	cdbHandle *cdb.CDB

	testRecordsStr []testCDBRecord
	testRecordsInt []testCDBRecord

	module *Module
}

func TestOCTestSuite(t *testing.T) {
	suite.Run(t, new(OCTestSuite))
}

func (suite *OCTestSuite) getCDBReader() cdb.Reader {
	// initialize reader
	reader, err := suite.cdbHandle.GetReader(suite.cdbFile)
	suite.Require().Nilf(err, "Can't get CDB reader: %#v", err)
	return reader
}

func (suite *OCTestSuite) getCDBWriter() cdb.Writer {
	// initialize writer
	writer, err := suite.cdbHandle.GetWriter(suite.cdbFile)
	suite.Require().Nilf(err, "Can't get CDB writer: %#v", err)
	return writer
}

func (suite *OCTestSuite) SetupTest() {
	f, err := ioutil.TempFile("", "test_*.cdb")
	suite.Require().Nilf(err, "Can't open temporary file: %#v", err)

	suite.cdbFile = f
	suite.cdbHandle = cdb.New() // create new cdb handle

	// generate test records
	tesRecsCnt := 2
	testRecordsStr := make([]testCDBRecord, tesRecsCnt)
	testRecordsInt := make([]testCDBRecord, tesRecsCnt)

	for i := 0; i < tesRecsCnt; i++ {
		stri := strconv.Itoa(i)
		typeByte := "s"
		testRecordsStr[i].key = []byte("/test/onlineconf/str" + stri)
		testRecordsStr[i].val = []byte(typeByte + "val" + stri)

		log.Printf("key %s val %s", string(testRecordsStr[i].key), string(testRecordsStr[i].val))

		testRecordsInt[i].key = []byte("/test/onlineconf/int" + stri)
		testRecordsInt[i].val = []byte(typeByte + stri)

		log.Printf("key %s val %s", string(testRecordsInt[i].key), string(testRecordsInt[i].val))

	}
	suite.testRecordsStr = testRecordsStr
	suite.testRecordsInt = testRecordsInt

	suite.fillTestCDB()

	module, err := NewModuleFromFile(f.Name())
	suite.Nilf(err, "Cant init onlineconf module!: %#v", err)
	suite.module = module
}

func (suite *OCTestSuite) TearDownTest() {
	err := suite.module.Close()
	suite.Nilf(err, "Can't close module: %#v", err)

	err = suite.cdbFile.Close()
	suite.Nilf(err, "Can't close cdb file: %#v", err)

	err = os.Remove(suite.cdbFile.Name())
	suite.Nilf(err, "Can't remove cdb file: %#v", err)
}

func (suite *OCTestSuite) fillTestCDB() {

	writer := suite.getCDBWriter()
	defer func() {
		err := writer.Close()
		suite.Require().Nilf(err, "Can't close cdb writer: %#v", err)
	}()

	allTestRecords := []testCDBRecord{}
	allTestRecords = append(allTestRecords, suite.testRecordsInt...)
	allTestRecords = append(allTestRecords, suite.testRecordsStr...)
	for _, rec := range allTestRecords {
		log.Printf("putting: key %s val %s", string(rec.key), string(rec.val))
		err := writer.Put(rec.key, rec.val)
		suite.Require().Nilf(err, "Cant put new value to cdb: %#v", err)
	}
}

func (suite *OCTestSuite) TestInt() {
	for _, testRec := range suite.testRecordsInt {
		ocInt, ok := suite.module.Int(string(testRec.key))
		suite.True(ok, "Cant find key %s in test onlineconf module %s", string(testRec.key), suite.module.Name)
		testInt, err := strconv.Atoi(string(testRec.val[1:]))
		if err != nil {
			panic(fmt.Errorf("Cant parse test record int: %w", err))
		}
		suite.Equal(ocInt, testInt)

		ocInt, ok = suite.module.IntWithDef(string(testRec.key), 0)
		suite.True(ok, "Cant find key %s in test onlineconf module %s", string(testRec.key), suite.module.Name)
		suite.Equal(ocInt, testInt)
	}

	for i, testRec := range suite.testRecordsInt {
		ocInt, ok := suite.module.IntWithDef(string(testRec.key)+"_not_exists", i)
		suite.False(ok, "Cant find key %s in test onlineconf module %s", string(testRec.key), suite.module.Name)
		suite.Equal(ocInt, i, "Default result was returned")
	}
}

func (suite *OCTestSuite) TestString() {
	for _, testRec := range suite.testRecordsStr {
		ocStr, ok := suite.module.String(string(testRec.key))
		suite.True(ok, "Cant find key %s in test onlineconf module %s", string(testRec.key), suite.module.Name)
		suite.Equal(string(testRec.val[1:]), ocStr)

		ocStr, ok = suite.module.StringWithDef(string(testRec.key), "")
		suite.True(ok, "Cant find key %s in test onlineconf module %s", string(testRec.key), suite.module.Name)
		suite.Equal(string(testRec.val[1:]), ocStr)

	}

	for i, testRec := range suite.testRecordsStr {
		defaultParamValue := "test_not_exists_" + strconv.Itoa(i)
		ocStr, ok := suite.module.StringWithDef(string(testRec.key)+"_not_exists", defaultParamValue)
		suite.False(ok, "Cant find key %s in test onlineconf module %s", string(testRec.key), suite.module.Name)
		suite.Equal(ocStr, defaultParamValue, "Default result was returned")
	}
}

func (suite *OCTestSuite) TestReload() {
	// todo
}

func (suite *OCTestSuite) TestUnknownParamType() {
	// todo
}

func (suite *OCTestSuite) TestJSON() {
	// todo
}

func (suite *OCTestSuite) TestConcurrent() {
	// todo
}
