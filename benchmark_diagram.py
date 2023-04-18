import matplotlib.pyplot as plt

# SingleKey
singleKey_data_size = [2e3, 2e4, 2e5]
singleKey_mutex = [12625, 207748, 205604]
singleKey_channel = [41798, 251148, 251540]
singleKey_syncMap = [15768, 66449, 68377]
singleKey_singleflight = [17039, 242412, 254424]

# 繪製 SingleKey 圖表
plt.plot(singleKey_data_size, singleKey_mutex, 'o-', label='Mutex')
plt.plot(singleKey_data_size, singleKey_channel, 'o-', label='Channel')
plt.plot(singleKey_data_size, singleKey_syncMap, 'o-', label='SyncMap')
plt.plot(singleKey_data_size, singleKey_singleflight, 'o-', label='Singleflight')
plt.xscale('log')
plt.title('Performance for SingleKey Scenario')
plt.xlabel('Data Size', fontsize=12)
plt.ylabel('ns/op', fontsize=12)
plt.legend()

# 設定 x 軸刻度
plt.xticks(singleKey_data_size, singleKey_data_size)
# 設定 x 軸刻度文字
plt.gca().set_xticklabels([f'{int(x):,}' for x in singleKey_data_size])

plt.tight_layout()
plt.show()

# MultiKey
multiKey_data_size = [2e3, 2e4, 2e5]
multiKey_mutex = [9026, 315541, 320547]
multiKey_channel = [19094, 16734, 33913]
multiKey_syncMap = [6556, 15004, 38783]
multiKey_singleflight = [11600, 22518, 98216]

# 繪製 MultiKey 圖表
plt.plot(multiKey_data_size, multiKey_mutex, 'o-', label='Mutex')
plt.plot(multiKey_data_size, multiKey_channel, 'o-', label='Channel')
plt.plot(multiKey_data_size, multiKey_syncMap, 'o-', label='SyncMap')
plt.plot(multiKey_data_size, multiKey_singleflight, 'o-', label='Singleflight')
plt.xscale('log')
plt.title('Performance for MultiKey Scenario')
plt.xlabel('Data Size', fontsize=12)
plt.ylabel('ns/op', fontsize=12)
plt.legend()

# 設定 x 軸刻度
plt.xticks(multiKey_data_size, multiKey_data_size)
# 設定 x 軸刻度文字
plt.gca().set_xticklabels([f'{int(x):,}' for x in multiKey_data_size])

plt.tight_layout()
plt.show()
