import { distributedKVStore, BusinessError } from '@kit.ArkData';
import { Logger } from '../utils/Logger';

let kvStore: distributedKVStore.SingleKVStore | null = null;

try {
  const kvStoreConfig: distributedKVStore.KVStoreConfig<distributedKVStore.Constants.INDEX_MODE> = {
    name: 'STORE_ID',
    type: distributedKVStore.KVStoreType.SINGLEKVSTORE,
    indexMode: distributedKVStore.IndexMode.INDEX_MODE_JSON,
    indexType: distributedKVStore.IndexType.INDEX_TYPE_STRING_KEY,
  };
  kvManager.getKVStore<distributedKVStore.SingleKVStore>(kvStoreConfig, async (err, store: distributedKVStore.SingleKVStore | undefined) => {
    if (err) {
      Logger.error(`Failed to get KVStore. Code:${err.code},message:${err.message}`);
      return;
    }
    Logger.info('Succeeded in getting KVStore.');
    kvStore = store;
  });
} catch (e) {
  let error = e as BusinessError;
  Logger.error(`An unexpected error occurred. Code:${error.code},message:${error.message}`);
}