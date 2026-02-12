import { browser } from '$app/environment';
import { initializeApp, type FirebaseApp, getApps, getApp } from 'firebase/app';
import {
    getAuth,
    signInWithPopup,
    signOut as firebaseSignOut,
    onAuthStateChanged,
    GoogleAuthProvider,
    FacebookAuthProvider,
    GithubAuthProvider,
    TwitterAuthProvider,
    type Auth,
    type User,
    type UserCredential,
    signInWithRedirect,
    getRedirectResult,
    createUserWithEmailAndPassword,
    signInWithEmailAndPassword,
    sendEmailVerification,
    sendPasswordResetEmail,
    confirmPasswordReset,
    updatePassword,
    updateProfile,
    EmailAuthProvider,
    reauthenticateWithCredential
} from 'firebase/auth';
import { firebaseClientConfig } from '@/constants/firebase.js';