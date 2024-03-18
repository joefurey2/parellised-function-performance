import pandas as pd
import matplotlib.pyplot as plt

# Create a DataFrame from the given dataset
data = {
    "Function type": ["OpenFaaS Sequential"] * 15,
    "Matrix Size": [100, 100, 100, 10000, 10000, 10000, 10000, 10000, 10000, 10000, 10000, 10000, 1000000, 1000000, 1000000],
    "Time elapsed": [0.0000241, 0.0000241, 0.0000256, 0.003301, 0.00395704, 0.0032646, 0.00437256, 0.0032953, 0.00320279, 0.00321569, 0.00314719, 0.00317989, 5.81728159, 5.36449185, 5.92909774]
}
df = pd.DataFrame(data)

# Plot the cluster chart
plt.figure(figsize=(10, 6))
plt.scatter(df["Matrix Size"], df["Time elapsed"], marker='o')
plt.title("Cluster Chart of Matrix Size by Time Elapsed")
plt.xlabel("Matrix Size")
plt.ylabel("Time Elapsed")
plt.xscale("log")  # Use logarithmic scale for better visualization
plt.yscale("log")
plt.grid(True)
plt.show()
