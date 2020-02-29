#CSVPARSER

This is a pet project fo a basic script to parse `.xml` files into csv files to import it to <img src="./assets/spendeelogo.jpeg" width="25"> [Spendee](https://www.spendee.com/).

#Usage

Add the column configuration json file to the `config` folder with the following name format `${bankId}Config.json`. The file should contain the column positions you whish to parse in a csv file from the xml source file i.e:

```
{
  "date": 0,
  "category": 2,
  "note": 3,
  "amount": 6
}
```

Download a bank `.xls` file in the `input` folder and run `./parse.sh`. the terminal should indicate how many opertions have been found and write them down in a file called `output.csv`


#Contributions

Any contributions are welcome, raise a PR if you want to improve the code.