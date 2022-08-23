'use strict';
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const bodyParser = require('body-parser');
const http = require('http')
const util = require('util');
const express = require('express')
const app = express();
const expressJWT = require('express-jwt');
const jwt = require('jsonwebtoken');
const bearerToken = require('express-bearer-token');
const cors = require('cors');
const constants = require('./config/constants.json')
const { 
    v1: uuidv1,
    v4: uuidv4,
} = require('uuid');

const host = process.env.HOST || constants.host;
const port = process.env.PORT || constants.port;


const helper = require('./app/helper')
const invoke = require('./app/invoke')
const qscc = require('./app/qscc')
const query = require('./app/query')

const mysql = require('mysql2'); 
const databaseCon = mysql.createConnection(
    {
        host: "localhost",
        user: "root",
        password: "root",
        database: "mydb"
    }
);

var mongo = require('mongodb');
var MongoClient = require('mongodb').MongoClient;
var url = "mongodb://localhost:27017/mydb";
MongoClient.connect(url, function(err, db) {
  if (err) throw err;
  console.log("Database created!");
  db.close();
});

app.options('*', cors());
app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: false
}));
// set secret variable
app.set('secret', 'thisismysecret');
app.use(expressJWT({
    secret: 'thisismysecret'
}).unless({
    path: ['/users','/users/login', '/register']
}));
app.use(bearerToken());

logger.level = 'debug';


app.use((req, res, next) => {
    logger.debug('New req for %s', req.originalUrl);
    if (req.originalUrl.indexOf('/users') >= 0 || req.originalUrl.indexOf('/users/login') >= 0 || req.originalUrl.indexOf('/register') >= 0) {
        return next();
    }
    var token = req.token;
    jwt.verify(token, app.get('secret'), (err, decoded) => {
        if (err) {
            console.log(`Error ================:${err}`)
            res.send({
                success: false,
                message: 'Failed to authenticate token. Make sure to include the ' +
                    'token returned from /users call in the authorization header ' +
                    ' as a Bearer token'
            });
            return;
        } else {
            req.username = decoded.username;
            req.orgname = decoded.orgName;
            logger.debug(util.format('Decoded from JWT token: username - %s, orgname - %s', decoded.username, decoded.orgName));
            return next();
        }
    });
});

var server = http.createServer(app).listen(port, function () { console.log(`Server started on ${port}`) });
logger.info('****************** SERVER STARTED ************************');
logger.info('***************  http://%s:%s  ******************', host, port);
server.timeout = 240000;

function getErrorMessage(field) {
    var response = {
        success: false,
        message: field + ' field is missing or Invalid in the request'
    };
    return response;
}

// Register and enroll user
app.post('/users', async function (req, res) {
    var username = req.body.username;
    var password = req.body.password;
    var orgName = req.body.orgName;
    logger.debug('End point : /users');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }

    databaseCon.connect(
        function(err){
            if (err) {
                console.log(err);
                return;
            }
            console.log("Connected to database.");
            var sql = "SELECT * FROM organizations WHERE `username` = '"+username+"' AND `password` = '"+password+"'";
            databaseCon.query(sql, async function (err, result) {
                var resLength = 0;
                if (err) {
                    console.log(err);
                    // return;
                } else {
                    resLength = result.length;
                }
                if (resLength != 0 || (username == "Admin@username" && password == "12345678")){
                    var token = jwt.sign({
                        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
                        username: username,
                        orgName: orgName
                    }, app.get('secret'));
                
                    let response = await helper.getRegisteredUser(username, orgName, true);
                
                    logger.debug('-- returned from registering the username %s for organization %s', username, orgName);
                    if (response && typeof response !== 'string') {
                        logger.debug('Successfully registered the username %s for organization %s', username, orgName);
                        response.token = token;
                        res.json(response);
                    } else {
                        logger.debug('Failed to register the username %s for organization %s with::%s', username, orgName, response);
                        res.json({ success: false, message: response });
                    }
                } else {
                    res.json({ success: true, message: "Invalid username or password."});
                }

            });
        }
    );
    

});

app.post('/database/query', async function(req, res){

    var sql = req.body.sql;
    
    databaseCon.connect(
        function(err){
            if (err) {
                console.log(err);
                return;
            }
            console.log("Connected!");
            //var sql = "CREATE TABLE organizations (id INT AUTO_INCREMENT PRIMARY KEY, username VARCHAR(255), password VARCHAR(255), role VARCHAR(255))";
            //var sql = "DROP TABLE organizations";
            databaseCon.query(sql, function (err, result) {
                if (err) {
                    console.log(err);
                    res.json({ success: false, message: err});
                    return;
                }
                res.json({ success: true, message: result});
            });
        }
    );

});

app.post('/organizations/insert', async function(req, res){

    var usernames = req.body.usernames;
    var passwords = req.body.passwords;
    var roles = req.body.roles;

    console.log(usernames);
    console.log(passwords);
    console.log(roles);

    if (usernames.length != passwords.length || usernames.length != roles.length || passwords.length != roles.length) {
        res.json({ success: false, message: "Invalid input arguments." });
        return;
    }

    databaseCon.connect(
        function(err){
            if (err) {
                console.log(err);
                return;
            }
            console.log("Connected to database.");
            var sql = "INSERT INTO organizations (username, password, role) VALUES ?";
            var values = [];
            for(var i = 0; i < usernames.length; i++) {
                values.push([usernames[i], passwords[i], roles[i]]);
            }
            databaseCon.query(sql,[values], function (err, result) {
                if (err) {
                    console.log(err);
                    return;
                }
                console.log("Number of records inserted: " + result.affectedRows);
                res.json({ success: true, message: "Number of records inserted: " + result.affectedRows});
            });
        }
    );

});

app.get('/organizations', async function(req, res){

    databaseCon.connect(
        function(err){
            if (err) {
                console.log(err);
                return;
            }
            console.log("Connected to database.");
            var sql = "SELECT * FROM organizations";
            databaseCon.query(sql, function (err, result) {
                if (err) {
                    console.log(err);
                    return;
                }
                res.json({ success: true, message: result});
            });
        }
    );

});

// Register and enroll user
app.post('/register', async function (req, res) {
    var username = req.body.username;
    var orgName = req.body.orgName;
    logger.debug('End point : /users');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }

    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
        username: username,
        orgName: orgName
    }, app.get('secret'));

    console.log(token)

    let response = await helper.registerAndGerSecret(username, orgName);

    logger.debug('-- returned from registering the username %s for organization %s', username, orgName);
    if (response && typeof response !== 'string') {
        logger.debug('Successfully registered the username %s for organization %s', username, orgName);
        response.token = token;
        res.json(response);
    } else {
        logger.debug('Failed to register the username %s for organization %s with::%s', username, orgName, response);
        res.json({ success: false, message: response });
    }

});

// Login and get jwt
app.post('/users/login', async function (req, res) {
    var username = req.body.username;
    var orgName = req.body.orgName;
    logger.debug('End point : /users');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }

    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
        username: username,
        orgName: orgName
    }, app.get('secret'));

    let isUserRegistered = await helper.isUserRegistered(username, orgName);

    if (isUserRegistered) {
        res.json({ success: true, message: { token: token } });

    } else {
        res.json({ success: false, message: `User with username ${username} is not registered with ${orgName}, Please register first.` });
    }
});

//----------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------
//-----------------------------Invoke transaction on smart contract-----------------------------
//----------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------

app.post('/channels/:channelName/chaincodes/:chaincodeName/chicken/owner/change', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var chickenId = req.body.chickenId;
        var newOwner = req.body.newOwner;

        let owner = req.username;

        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('chickenId  : ' + chickenId);
        logger.debug('newOwner  : ' + newOwner);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!chickenId) {
            res.json(getErrorMessage('\'chickenId\''));
            return;
        }
        if (!newOwner) {
            res.json(getErrorMessage('\'newOwner\''));
            return;
        }

        let message = await invoke.changeChickenOwner(channelName, chaincodeName, req.username, req.orgname, chickenId, owner, newOwner);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.post('/channels/:channelName/chaincodes/:chaincodeName/chicken/create', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var birthday = req.body.birthday;
        var breed = req.body.breed;
        var price = req.body.price;

        let owner = req.username;
        let id = uuidv4();

        console.log("id========================"+id)

        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('birthday  : ' + birthday);
        logger.debug('breed  : ' + breed);
        logger.debug('price  : ' + price);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!birthday) {
            res.json(getErrorMessage('\'birthday\''));
            return;
        }
        if (!breed) {
            res.json(getErrorMessage('\'breed\''));
            return;
        }
        if (!price) {
            res.json(getErrorMessage('\'price\''));
            return;
        }


        let message = await invoke.createChicken(channelName, chaincodeName, req.username, req.orgname, id, birthday, breed, price, owner);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.post('/channels/:channelName/chaincodes/:chaincodeName/token/buy', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var price = req.body.price;

        let user = req.username;

        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('price  : ' + price);
        logger.debug('user  : ' + user);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!price) {
            res.json(getErrorMessage('\'price\''));
            return;
        }

        let message = await invoke.buyToken(channelName, chaincodeName, req.username, req.orgname, user, price);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.post('/channels/:channelName/chaincodes/:chaincodeName/chicken/public', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var assetId = req.body.assetId;

        let user = req.username;

        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('assetId  : ' + assetId);
        logger.debug('user  : ' + user);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!assetId) {
            res.json(getErrorMessage('\'assetId\''));
            return;
        }

        let message = await invoke.setChickenPublicForBuy(channelName, chaincodeName, req.username, req.orgname, assetId, user);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.post('/channels/:channelName/chaincodes/:chaincodeName/token/transfer', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var amount = req.body.amount;
        var receiver = req.body.receiver;

        let user = req.username;

        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('amount  : ' + amount);
        logger.debug('receiver  : ' + receiver);
        logger.debug('user  : ' + user);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!amount) {
            res.json(getErrorMessage('\'amount\''));
            return;
        }
        if (!receiver) {
            res.json(getErrorMessage('\'receiver\''));
            return;
        }

        let message = await invoke.transferToken(channelName, chaincodeName, req.username, req.orgname, user, receiver, amount);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.post('/channels/:channelName/chaincodes/:chaincodeName/asset/bid', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var assetId = req.body.assetId;
        var assetOwner = req.body.assetOwner;
        var price = req.body.price;

        let customer = req.username;

        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('assetId  : ' + assetId);
        logger.debug('assetOwner  : ' + assetOwner);
        logger.debug('customer  : ' + customer);
        logger.debug('price  : ' + price);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!assetId) {
            res.json(getErrorMessage('\'assetId\''));
            return;
        }
        if (!assetOwner) {
            res.json(getErrorMessage('\'assetOwner\''));
            return;
        }
        if (!price) {
            res.json(getErrorMessage('\'price\''));
            return;
        }

        let message = await invoke.bidForAsset(channelName, chaincodeName, req.username, req.orgname, assetId , customer , assetOwner , price);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.post('/channels/:channelName/chaincodes/:chaincodeName/chicken/sell', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var id = req.body.id;
        var customer = req.body.customer;

        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('id  : ' + id);
        logger.debug('customer  : ' + customer);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!id) {
            res.json(getErrorMessage('\'id\''));
            return;
        }
        if (!customer) {
            res.json(getErrorMessage('\'customer\''));
            return;
        }

        let message = await invoke.sellChicken(channelName, chaincodeName, req.username, req.orgname, id, req.username, customer);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.post('/channels/:channelName/chaincodes/:chaincodeName/chicken/price', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var id = req.body.id;
        var price = req.body.price;

        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('id  : ' + id);
        logger.debug('price  : ' + price);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!id) {
            res.json(getErrorMessage('\'id\''));
            return;
        }
        if (!price) {
            res.json(getErrorMessage('\'price\''));
            return;
        }

        let message = await invoke.setChickenPrice(channelName, chaincodeName, req.username, req.orgname, id, price);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.post('/channels/:channelName/chaincodes/:chaincodeName/information/growth', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var id = req.body.id;
        var key = req.body.key;
        var value = req.body.value;
        var instruction = req.body.instruction;

        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('id  : ' + id);
        logger.debug('key  : ' + key);
        logger.debug('value  : ' + value);
        logger.debug('instruction  : ' + instruction);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!id) {
            res.json(getErrorMessage('\'id\''));
            return;
        }
        if (!key) {
            res.json(getErrorMessage('\'key\''));
            return;
        }
        if (!value) {
            res.json(getErrorMessage('\'value\''));
            return;
        }
        if (!instruction) {
            res.json(getErrorMessage('\'instruction\''));
            return;
        }

        let message = await invoke.putGrowthInformation(channelName, chaincodeName, req.username, req.orgname, id, key, value, instruction);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});


//----------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------
//----------------------------------Query on smart contract-------------------------------------
//----------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------

app.get('/channels/:channelName/chaincodes/:chaincodeName/chickens/owner', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`);

        let owner = req.username;

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        
        let message = await query.queryChickensByOwner(channelName, chaincodeName, req.username, req.orgname, owner);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/channels/:channelName/chaincodes/:chaincodeName/chicken', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`);
        var chickenId = req.query.chickenId;
        console.log(`chickenId is :${chickenId}`);

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!chickenId) {
            res.json(getErrorMessage('\'chickenId\''));
            return;
        }


        let message = await query.queryChicken(channelName, chaincodeName, req.username, req.orgname, chickenId);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/channels/:channelName/chaincodes/:chaincodeName/asset/history', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`);
        var chickenId = req.query.chickenId;

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!chickenId) {
            res.json(getErrorMessage('\'chickenId\''));
            return;
        }

        let message = await query.getHistoryForChicken(channelName, chaincodeName, req.username, req.orgname, chickenId);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/channels/:channelName/chaincodes/:chaincodeName/chickens/all', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`);

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }

        let message = await query.queryAllChickens(channelName, chaincodeName, req.username, req.orgname);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/channels/:channelName/chaincodes/:chaincodeName/token', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`);

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }

        let message = await query.queryToken(channelName, chaincodeName, req.username, req.orgname, req.username);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/channels/:channelName/chaincodes/:chaincodeName/asset/bids', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        var assetId = req.query.assetId;
        console.log(`chaincode name is :${chaincodeName}`);

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }

        let message = await query.queryBidsOfAsset(channelName, chaincodeName, req.username, req.orgname, assetId, req.username);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/channels/:channelName/chaincodes/:chaincodeName/chickens/public', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;

        console.log(`chaincode name is :${chaincodeName}`);

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }

        let message = await query.queryPublicChickens(channelName, chaincodeName, req.username, req.orgname);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/qscc/channels/:channelName/chaincodes/:chaincodeName', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`)
        let args = req.query.args;
        let fcn = req.query.fcn;
        // let peer = req.query.peer;

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('fcn : ' + fcn);
        logger.debug('args : ' + args);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!fcn) {
            res.json(getErrorMessage('\'fcn\''));
            return;
        }
        if (!args) {
            res.json(getErrorMessage('\'args\''));
            return;
        }
        console.log('args==========', args);
        args = args.replace(/'/g, '"');
        args = JSON.parse(args);
        logger.debug(args);

        let response_payload = await qscc.qscc(channelName, chaincodeName, args, fcn, req.username, req.orgname);

        // const response_payload = {
        //     result: message,
        //     error: null,
        //     errorData: null
        // }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});