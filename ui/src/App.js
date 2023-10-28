import React from "react";
import logo from './logo.svg';
import './App.css';


class App extends React.Component {
    state = {
        total: null,
        next: null,
        operation: null,
    };

    handleClick = buttonName => {
        var xmlhttp = new XMLHttpRequest();
        xmlhttp.open('POST', 'http://localhost:8088/axis2/services/test-1/', true);

        // build SOAP request
        var sr =
            '<?xml version="1.0" encoding="utf-8"?>' +
            '<soapenv:Envelope ' +
            'xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ' +
            'xmlns:api="http://127.0.0.1/Integrics/Enswitch/API" ' +
            'xmlns:xsd="http://www.w3.org/2001/XMLSchema" ' +
            'xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">' +
            '<soapenv:Body>' +
            '<api:some_api_call soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">' +
            '<username xsi:type="xsd:string">login_username</username>' +
            '<password xsi:type="xsd:string">password</password>' +
            '</api:some_api_call>' +
            '</soapenv:Body>' +
            '</soapenv:Envelope>';

        xmlhttp.onreadystatechange = function () {
            if (xmlhttp.readyState === 4) {
                if (xmlhttp.status === 200) {
                    alert(xmlhttp.responseText);
                    // alert('done. use firebug/console to see network response');
                }
            }
        }
        // Send the POST request
        xmlhttp.setRequestHeader('Content-Type', 'text/xml');
        xmlhttp.send(sr);
    };

    render() {
        return (
            <div className="App">
                <header className="App-header">
                    <img src={logo} className="App-logo" alt="logo"/>
                    <p>
                        Edit <code>src/App.js</code> and save to reload.
                    </p>
                    <a
                        className="App-link"
                        href="https://reactjs.org"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        Learn React
                    </a>
                    <div onClick={this.handleClick}>123123123</div>
                </header>
            </div>
        );
    }
}

export default App;
