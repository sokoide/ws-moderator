"use client";
import { createContext, useEffect, useState, ReactNode} from "react";

interface AppContextType {
    loginInfo?: {
        loggedIn: boolean;
        email: string;
    };
    login: (email: string) => void;
    logout: () => void;
}
export const AppContext = createContext<AppContextType | null>(null);

interface AppContextProviderProps {
    children: ReactNode;
}

// functions
const getLoginInfo = () => {
    let loginInfo = {
        loggedIn: false,
        email: "",
    };

    if (typeof window !== "undefined") {
        const stored = localStorage.getItem("loginInfo");
        if (stored == null) return loginInfo;

        try {
            let ret = JSON.parse(stored);
            return ret;
        } catch (e) {
            return loginInfo;
        }
    }
    console.log(
        "getLoginInfo: window NOT available, using the empty login info"
    );
    return loginInfo;
};

// ContextProvider
export const AppContextProvider: React.FC<AppContextProviderProps> = ({ children }) => {
    // ----- states -----
    const [loginInfo, setLoginInfo] = useState(getLoginInfo());

    // save loginInfo in localStorage whenever it changes
    useEffect(() => {
        console.log("saving loginInfo %O in localStorage", loginInfo);
        localStorage.setItem("loginInfo", JSON.stringify(loginInfo));
    }, [loginInfo]);

    // --- functions ---
    const login = (email: string) => {
        console.log("login(%O)", email);

        setLoginInfo((prev: AppContextType["loginInfo"] | undefined) => ({
            ...prev,
            loggedIn: true,
            email: email,
        }));
    };

    const logout = () => {
        console.log("logout");
        setLoginInfo((prev: AppContextType["loginInfo"] | undefined) => ({
            ...prev,
            loggedIn: false,
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
            {children}
        </AppContext.Provider>
    );
};
