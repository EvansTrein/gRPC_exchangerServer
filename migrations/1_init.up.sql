CREATE TABLE IF NOT EXISTS Currencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  
    code VARCHAR(5) NOT NULL UNIQUE,             
    name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS Rates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    baseCurrencyID INTEGER NOT NULL,     
    toCurrencyID INTEGER NOT NULL,   
    rate FLOAT NOT NULL,                    
    date DATE NOT NULL,                    
    FOREIGN KEY (baseCurrencyID) REFERENCES Currencies(id),
    FOREIGN KEY (toCurrencyID) REFERENCES Currencies(id),
    UNIQUE (baseCurrencyID, toCurrencyID, date)  
);