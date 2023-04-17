import matplotlib.pyplot as plt

# OneKey
onekey_data_size = [2e3, 2e4, 2e5]
onekey_mutex = [12625, 207748, 205604]
onekey_channel = [41798, 251148, 251540]
onekey_syncmap = [15768, 66449, 68377]
onekey_singleflight = [17039, 242412, 254424]

# MultiKey
multikey_data_size = [2e3, 2e4, 2e5]
multikey_mutex = [9026, 315541, 320547]
multikey_channel = [19094, 16734, 33913]
multikey_syncmap = [6556, 15004, 38783]
multikey_singleflight = [11600, 22518, 98216]

# 繪製 OneKey 圖表
plt.plot(onekey_data_size, onekey_mutex, 'o-', label='Mutex')
plt.plot(onekey_data_size, onekey_channel, 'o-', label='Channel')
plt.plot(onekey_data_size, onekey_syncmap, 'o-', label='SyncMap')
plt.plot(onekey_data_size, onekey_singleflight, 'o-', label='Singleflight')
plt.xscale('log')
plt.title('OneKey Performance')
plt.xlabel('Data Size', fontsize=12)
plt.ylabel('ns/op', fontsize=12)
plt.legend()

# 設定 x 軸刻度
plt.xticks(onekey_data_size, onekey_data_size)
# 設定 x 軸刻度文字
plt.gca().set_xticklabels([f'{int(x):,}' for x in onekey_data_size])

plt.tight_layout()
plt.show()

# 繪製 MultiKey 圖表
plt.plot(multikey_data_size, multikey_mutex, 'o-', label='Mutex')
plt.plot(multikey_data_size, multikey_channel, 'o-', label='Channel')
plt.plot(multikey_data_size, multikey_syncmap, 'o-', label='SyncMap')
plt.plot(multikey_data_size, multikey_singleflight, 'o-', label='Singleflight')
plt.xscale('log')
plt.title('MultiKey Performance')
plt.xlabel('Data Size', fontsize=12)
plt.ylabel('ns/op', fontsize=12)
plt.legend()

# 設定 x 軸刻度
plt.xticks(multikey_data_size, multikey_data_size)
# 設定 x 軸刻度文字
plt.gca().set_xticklabels([f'{int(x):,}' for x in multikey_data_size])

plt.tight_layout()
plt.show()
