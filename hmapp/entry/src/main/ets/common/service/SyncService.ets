import { Logger } from '../../utils/Logger'
import { distributedKVStore, BusinessError } from '@kit.ArkData';

let kvStore: distributedKVStore.SingleKVStore | null = null;
let kvManager: distributedKVStore.KVManager | null = null;
const BASE_KEY_COUNT = 1000;

try {
  const kvManagerConfig: distributedKVStore.KVManagerConfig = {
    context: this.context,
    bundleName: 'com.example.app'
  };
  kvManager = distributedKVStore.createKVManager(kvManagerConfig);
  let kvStoreConfig: distributedKVStore.KVStoreConfig<distributedKVStore.Constants.INDEX_MODE> = {
    name: 'STORE_ID',
    type: distributedKVStore.KVStoreType.SINGLEKVSTORE
  };
  kvManager!.getKVStore<distributedKVStore.SingleKVStore>(kvStoreConfig, (err, store: distributedKVStore.SingleKVStore) => {
    if (err != undefined) {
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

let arr = new Uint8Array([21, 31]);
let valueBlob: distributedKVStore.ValueType = arr;

try {
  kvStore!.put(BASE_KEY_COUNT++, valueBlob, (err) => {
    if (err != undefined) {
      Logger.error(`Failed to put Blob. Code:${err.code},message:${err.message}`);
      return;
    }
    Logger.info('Succeeded in putting Blob');
  });
} catch (e) {
  let error = e as BusinessError;
  Logger.error(`An unexpected error occurred. Code:${error.code},message:${error.message}`);
}