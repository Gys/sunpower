SunPower
========

Calculate the average insolation (exposure to the sun) for any location in the world.

**Used Data**

The data is retrieved from NASA POWER Project Data Sets. Where POWER stands *Prediction Of Worldwide Energy Resources*. It has solar and meteorological data sets for support of renewable energy, building energy efficiency and agricultural needs.

See: https://power.larc.nasa.gov/docs/v1/

The NASA API provides a lot of data. In this case the ALLSKY_SFC_SW_DWN is used. It means *All Sky Insolation Incident on a Horizontal Surface*. A value expressed in kW-hr/m^2/day.

**Why**

Sau Sheong Chang created a detailed blog post on the energy requirement of Singapore and how installing solar panels on all rooftops could help. 
https://medium.com/sausheong/estimate-the-solar-output-of-your-rooftop-with-google-maps-725e4f636f14
He calculated the average kW-hr/m^2/day by hand and as I was wondering what the value could be for other locations, I decided to just create this small tool for that.

**Application**

Given a latitude and longitude SunPower will calculate the average kW-hr/m^2/day for that location, where the average is based on all daily values over a period of the past ten years.

Install:

    go get github/Gys/sunpower

Usage:

    ./sunpower 38.722501 -9.4323331


Same sample locations, calculated on Feb 7, 2020:

location| latitude longitude |kW-hr/m^2/day|
|---|---|---|
|Singapore | 1.352083 103.819839|4.55|
|Lisbon|38.722501 -9.4323331|4.80|
|Vancouver|49.282730 -123.120735|3.11|
|Miami|25.761681 -80.191788|5.01|
|Abu Dhabi|24.453884 54.377342|5.92|
|Stockholm|59.326242 17.8419719|2.62|


