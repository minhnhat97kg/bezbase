import { useState, useEffect } from 'react';
import { rbacService } from '../services/api';

interface ResourceActionOptions {
  resources: any[];
  actions: any[];
  loading: boolean;
}

export default function useResourceActionOptions(): ResourceActionOptions {
  const [resources, setResources] = useState<any[]>([]);
  const [actions, setActions] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

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
