import { useState, useEffect } from 'react';
import { rbacService } from '../services/api';

export default function useResourceActionOptions() {
  const [resources, setResources] = useState([]);
  const [actions, setActions] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchOptions() {
      setLoading(true);
      try {
        const resRes = await rbacService.getResources({ page_size: 100 });
        const actRes = await rbacService.getActions({ page_size: 100 });
        setResources(resRes.data.data || []);
        setActions(actRes.data.data || []);
      } catch (err) {
        setResources([]);
        setActions([]);
      } finally {
        setLoading(false);
      }
    }
    fetchOptions();
  }, []);

  return { resources, actions, loading };
}
