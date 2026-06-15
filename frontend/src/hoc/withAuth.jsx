import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import useAuthStore from '../store/authStore';
import { getToken } from '../utils/cookies';
import { getCurrentUser } from '../services/api';

function withAuth(WrappedComponent, { roles = [] } = {}) {
  function AuthGuard(props) {
    const { isAuthed, isLoading, user, login, stopLoading } = useAuthStore();
    const navigate = useNavigate();

    useEffect(() => {
      const initAuth = async () => {
        const token = await getToken();

        if (!token) {
          useAuthStore.getState().logout();
          stopLoading();
          return;
        }

        try {
          const res = await getCurrentUser();
          await login(token, res.data);
        } catch (error) {
          useAuthStore.getState().logout();
        } finally {
          stopLoading();
        }
      };

      initAuth();
    }, []);

    useEffect(() => {
      if (isLoading) return;

      if (!isAuthed) {
        navigate('/login', { replace: true });
        return;
      }

      if (roles.length > 0 && !roles.includes(user?.role)) {
        navigate('/dashboard', { replace: true });
      }
    }, [isAuthed, isLoading, user, navigate]);

    if (isLoading) return null;
    if (!isAuthed) return null;
    if (roles.length > 0 && !roles.includes(user?.role)) return null;

    return <WrappedComponent {...props} />;
  }

  AuthGuard.displayName = `withAuth(${WrappedComponent.displayName ?? WrappedComponent.name ?? 'Component'})`;
  return AuthGuard;
}

export default withAuth;
