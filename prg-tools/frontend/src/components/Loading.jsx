

function Loading({message = "Loading..."}) {

    return (
        <div className="flex flex-col items-center justify-center gap-3 py-8">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-blue-500" />
            <p className="text-sm text-gray-500">{message}</p>
        </div>
    )
}

export default Loading;