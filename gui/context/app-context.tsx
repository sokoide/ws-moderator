"use client";
import { createContext, useEffect, useState } from "react";

export const AppContext = createContext(null);

// functions
const getLoginInfo = () => {
    let loginInfo = {};

    if (typeof window !== "undefined") {
        const stored = localStorage.getItem("loginInfo") || {loggedIn: false, username: "", employee: "", email: ""};
        return stored ? JSON.parse(stored) : loginInfo;
    }
    console.log(
        "getLoginInfo: window NOT available, using the empty login info"
    );
    return loginInfo;
};

// ContextProvider
export const AppContextProvider = (props) => {
    // ----- states -----
    const [loginInfo, setLoginInfo] = useState(getLoginInfo());
    // save loginInfo in localStorage whenever it changes
    useEffect(() => {
        console.log("saving loginInfo %O in localStorage", loginInfo);
        localStorage.setItem("loginInfo", JSON.stringify(loginInfo));
    }, [loginInfo]);

    // --- functions ---
    const login = (username: string, employee: string, email: string) => {
        console.log("login(%O, %O, %O)", username, employee, email);

        setLoginInfo((prev) => ({
            ...prev,
            loggedIn: true,
            username: username,
            employee: employee,
            email: email,
        }));
    };

    const logout = () => {
        console.log("logout");
        setLoginInfo((prev) => ({
            ...prev,
            loggedIn: false,
            username: "",
            employee: "",
            email: "",
        }));
    };

    // contextValue
    const contextValue = {
        loginInfo,
        login,
        logout,
    };
    return (
        <AppContext.Provider value={contextValue}>
            {props.children}
        </AppContext.Provider>
    );
};
