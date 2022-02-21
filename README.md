## Calorie Insights

### Intro

* CLI to gain insights from a calorie counter CSV file. What kind of insights?

* Get weekly average calories consumed 
* Get percantage of calories contributed by specific sources (TO BE IMPLEMENTED)

* The file that's used for insights must be of this format:

```
01/09/22 7:11 PM,Pasta with veggies,400
01/09/22 7:11 PM,Ragi chips,420
01/09/22 7:11 PM,Homemade meal,520
```

* This is a free app that exports calorie data in this format: https://play.google.com/store/apps/details?id=com.doubleblacksoftware.caloriecounter&hl=en_US&gl=US.

* Using 'calorie-insights' alongside the above app can give interesting insights into calories consumed.

### Usage

* Build the binary:

```
go build -o calorie-insights

```

* Run the binary to view CLI options

```
./calorie-insights 

OR 

./calorie-insights --help
```
