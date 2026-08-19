import { useCallback, useEffect, useState } from 'react';
import { AuthService } from '../../bindings/github.com/songwei.ma/talus-mofish/backend/services';
import { UserProfile } from '../utils/userProfile';

type OAuthProvider = 'github' | 'google';
type AuthProvider = 'email' | OAuthProvider;

export function useCurrentUser() {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [signingIn, setSigningIn] = useState<AuthProvider | null>(null);

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const profile = (await AuthService.GetCurrentUser()) as UserProfile | null;
      setUser(profile ?? null);
    } catch (err) {
      console.error(err);
      setUser(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const signInWithEmail = useCallback(async (email: string) => {
    setSigningIn('email');
    try {
      const profile = (await AuthService.SignInWithEmail(email)) as UserProfile;
      setUser(profile);
      return profile;
    } finally {
      setSigningIn(null);
    }
  }, []);

  const signIn = useCallback(async (provider: OAuthProvider) => {
    setSigningIn(provider);
    try {
      const profile =
        provider === 'github'
          ? ((await AuthService.SignInWithGitHub()) as UserProfile)
          : ((await AuthService.SignInWithGoogle()) as UserProfile);
      setUser(profile);
      return profile;
    } finally {
      setSigningIn(null);
    }
  }, []);

  const signOut = useCallback(async () => {
    await AuthService.SignOut();
    setUser(null);
  }, []);

  return {
    user,
    loading,
    signingIn,
    refresh,
    signInWithEmail,
    signIn,
    signOut,
  };
}
